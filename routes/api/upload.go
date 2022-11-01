package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	Tools "github.com/NotRoyadma/auto_backup-dnxrg/avails"
	Conf "github.com/NotRoyadma/auto_backup-dnxrg/config"
	"github.com/NotRoyadma/auto_backup-dnxrg/logger"
	Static "github.com/NotRoyadma/auto_backup-dnxrg/routes/static"
	websocket "github.com/gorilla/websocket"
)

func UploadFileManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("token") != Conf.Conf.Token {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var nodeName string = r.Header.Get("node")
	if nodeName == "" {
		http.Error(w, "Node name must be specified", http.StatusBadRequest)
		return
	}

	if !Tools.StringInSlice(nodeName, Conf.Conf.Nodes) {
		http.Error(w, "This node is not specified. maybe because misconfig", http.StatusBadRequest)
		return
	}

	// All the validations are done
	// now we have `node` eg:node1,game2,in2 and we have `token`
	// since this request is to just get the file into data folder so we'll just check for data folder, the node's folder and its data.zip

	if DoDirExists, _ := Tools.DoDirectoryExists(Conf.Conf.DataDirectory); !DoDirExists {
		err := os.Mkdir(Conf.Conf.DataDirectory, 0777)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.WriteERRLog("api/upload.go 48 " + err.Error())
			return
		}
	}

	if DoDirExists, _ := Tools.DoDirectoryExists(Conf.Conf.DataDirectory + "/" + nodeName); !DoDirExists {
		err := os.Mkdir(Conf.Conf.DataDirectory+"/"+nodeName, 0777)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.WriteERRLog("api/upload.go 57 " + err.Error())
			return
		}
	}

	err := os.RemoveAll(Conf.Conf.DataDirectory + "/" + nodeName + "/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WriteERRLog("api/upload.go 65 " + err.Error())
		return
	}

	// ok so till here we have data folder and the node's folder
	// we also have deleted all the server files of the node
	// now lets upgrade the request to websockets and initialize the file transfer
	var webSocketUpgrader = websocket.Upgrader{}

	//upgrade the request
	ws, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WriteERRLog("api/upload.go 78 " + err.Error())
		return
	}

	// now we have got `ws` means socket with the help of gorilla websockets
	// now we will wait for the client to initiate the file transfer
	defer ws.Close()

	var InitialResponseFromWebsocket interface{}
	// waiting for new message
	// if after 8 minutes, new message don't arrives then we will close it
	for {
		ws.SetReadDeadline(time.Now().Add(8 * time.Minute))
		_, strdata, err := ws.ReadMessage()
		if err != nil {
			return
		}
		err = json.Unmarshal(strdata, &InitialResponseFromWebsocket)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Invaid JSON"))
			ws.Close()
			return
		}
		break
	}

	// got the client initialization message
	InitialResponseFromWebsocketJson, _ := InitialResponseFromWebsocket.(map[string]interface{})
	if InitialResponseFromWebsocketJson["event"] != "initiate_file" {
		logger.WriteERRLog("api/upload.go 107 " + InitialResponseFromWebsocket.(string))
		ws.Close()
		return
	}

	// initialization of file sending
	// now the file sending hasbeen started
	// the client is waiting for server's response to start sending subfolders
	type ServerResponse struct {
		Event string
		Chunk int
	}
	serverResponse, err := json.Marshal(ServerResponse{
		Event: "initiate_subfolders",
		Chunk: Conf.Conf.ChunkSize,
	})
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		ws.Close()
		return
	}

	ws.WriteMessage(websocket.TextMessage, serverResponse)

	// assuming that client got the response
	// now client will start looping through all the subfolders and send data as followed
	// Event: "subfolder_start"
	// Name: "folder name eg.Server UUID"
	// and then the client will wait for server
	// now we have to wait for client
	// looping through all the subfolders can take time upto max 20 minutes

	// we have to loop through all the subfolders the client is sending
	for {
		var subFolderData interface{}
		for {
			ws.SetReadDeadline(time.Now().Add(20 * time.Minute))
			_, strdata, err := ws.ReadMessage()
			if err != nil {
				return
			}
			err = json.Unmarshal(strdata, &subFolderData)
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("Invaid JSON"))
				ws.Close()
				return
			}
			break
		}

		// we have got the event subfolder_data as followed
		subFolderDataJson, _ := subFolderData.(map[string]interface{})

		//check if the client send's EOF
		if subFolderDataJson["event"] == "end_sharing" {
			break
		}

		if subFolderDataJson["event"] != "subfolder_start" {
			logger.WriteERRLog("api/upload.go 166 " + subFolderData.(string))
			ws.Close()
			return
		}

		FolderNameFromClient := subFolderDataJson["name"].(string)

		// create the folder at this endpoint
		err := os.Mkdir(Conf.Conf.DataDirectory+"/"+nodeName+"/"+FolderNameFromClient, 0777)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			ws.Close()
			logger.WriteERRLog("api/upload.go 178 " + err.Error())
			return
		}

		// till here the client is waiting for us to initiate to send the folder's data or if the folder is empy then the client will send event "folder_empy".

		type ServerResponse struct {
			Event      string
			FolderName string
		}

		serverResponse, err := json.Marshal(ServerResponse{
			Event:      "subfolder_data_start",
			FolderName: FolderNameFromClient,
		})

		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			ws.Close()
			return
		}

		// open the file to write
		WriteAbleFile, _ := os.OpenFile(Conf.Conf.DataDirectory+"/"+nodeName+"/"+FolderNameFromClient+"/"+"data.zip", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

		ws.WriteMessage(websocket.TextMessage, serverResponse)

		// now the client got the response to start sending data
		// it will take the data.zip and read the defined chunks and send it to server
		// now we have to wait for the client in a loop
		for {
			var subFolderChunkData interface{}

			// wait for client to send data in chunks

			for {
				ws.SetReadDeadline(time.Now().Add(20 * time.Minute))
				_, strdata, err := ws.ReadMessage()
				if err != nil {
					return
				}
				err = json.Unmarshal(strdata, &subFolderChunkData)
				if err != nil {
					ws.WriteMessage(websocket.TextMessage, []byte("Invaid JSON"))
					ws.Close()
					return
				}
				break
			}

			subFolderChunkDataJson := subFolderChunkData.(map[string]interface{})

			// handle if the chunk ended
			if subFolderChunkDataJson["event"] == "end_s_chunk" {
				break
			}

			if subFolderChunkDataJson["event"] != "subfolder_chunk_data" {
				logger.WriteERRLog("api/upload.go 236 " + subFolderChunkData.(string))
				ws.Close()
				return
			}

			// now we have got just a chunk of data and we to to write it
			parsedChunkFromRequest, _ := base64.StdEncoding.DecodeString(subFolderChunkDataJson["Chunk"].(string))

			// write the data to file :0
			WriteAbleFile.Write(parsedChunkFromRequest)

			// the client was waiting for server to send back the response after writing
			serverResponse, _ := json.Marshal(ServerResponse{
				Event:      "subfolder_chunk_data_ack",
				FolderName: FolderNameFromClient,
			})

			ws.WriteMessage(websocket.TextMessage, serverResponse)

		}

	}

	//the client

}
