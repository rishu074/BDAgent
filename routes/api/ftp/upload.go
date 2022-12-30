package ftp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	Tools "github.com/NotRoyadma/BDAgent/avails"
	Conf "github.com/NotRoyadma/BDAgent/config"
	"github.com/NotRoyadma/BDAgent/logger"
	Static "github.com/NotRoyadma/BDAgent/routes/static"
	humanize "github.com/dustin/go-humanize"
	websocket "github.com/gorilla/websocket"
	ftp "github.com/jlaffaye/ftp"
)

func UploadFileManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		logger.WriteERRLog(" /api/ftp/upload.go 21 Method now allowed")
		return
	}

	if r.Header.Get("token") != Conf.Conf.Token {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		logger.WriteERRLog(" /api/ftp/upload.go 27 Token not valid")
		return
	}

	var nodeName string = r.Header.Get("node")
	if nodeName == "" {
		http.Error(w, "Node name must be specified", http.StatusBadRequest)
		logger.WriteERRLog(" /api/ftp/upload.go 34 Node name must be specified")
		return
	}

	if !Tools.StringInSlice(nodeName, Conf.Conf.Nodes) {
		http.Error(w, "This node is not specified. maybe because misconfig", http.StatusBadRequest)
		logger.WriteERRLog(" /api/ftp/upload.go 40 No node with " + nodeName)
		return
	}

	// All the validations are done
	// now we have `node` eg:node1,game2,in2 and we have `token`
	// since this is the FTP one so we will connect to ftp first
	ftpClient, err := ConnectFtp()
	defer ftpClient.Quit()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.WriteERRLog(" /api/ftp/upload.go 52 " + err.Error())
		return
	}

	if DoDirExists, _ := Tools.DoDirExistsFTP(ftpClient, "/", Conf.Conf.DataDirectory); !DoDirExists {
		err := ftpClient.MakeDir(Conf.Conf.DataDirectory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.WriteERRLog("api/ftp/upload.go 60 " + err.Error())
			return
		}
	}

	if DoDirExists, _ := Tools.DoDirExistsFTP(ftpClient, "/"+Conf.Conf.DataDirectory, nodeName); !DoDirExists {
		err := ftpClient.MakeDir("/" + Conf.Conf.DataDirectory + "/" + nodeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.WriteERRLog("api/ftp/upload.go 69 " + err.Error())
			return
		}
	} else {
		err := ftpClient.RemoveDirRecur("/" + Conf.Conf.DataDirectory + "/" + nodeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.WriteERRLog("api/ftp/upload.go 76 " + err.Error())
			return
		}

		err = ftpClient.MakeDir("/" + Conf.Conf.DataDirectory + "/" + nodeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.WriteERRLog("api/ftp/upload.go 83 " + err.Error())
			return
		}
	}

	_ = ftpClient.ChangeDirToParent()

	// ok so till here we have data folder and the node's folder
	// we also have deleted all the server files of the node
	// now lets upgrade the request to websockets and initialize the file transfer
	var webSocketUpgrader = websocket.Upgrader{}

	//upgrade the request
	ws, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WriteERRLog("api/ftp/upload.go 99 " + err.Error())
		return
	}

	// now we have got `ws` means socket with the help of gorilla websockets
	// now we will wait for the client to initiate the file transfer
	var TotalBytes int64 = 0
	defer ws.Close()
	defer func() {
		s := humanize.Bytes(uint64(TotalBytes))
		logger.WriteLog(nodeName + " uploaded " + s)
	}()

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
	if InitialResponseFromWebsocketJson["Event"] != "initiate_file" {
		logger.WriteERRLog("api/ftp/upload.go 133 " + InitialResponseFromWebsocketJson["Event"].(string))
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
		if subFolderDataJson["Event"] == "end_sharing" {
			break
		}

		if subFolderDataJson["Event"] != "subfolder_start" {
			logger.WriteERRLog("api/ftp/upload.go 192 " + subFolderData.(string))
			ws.Close()
			return
		}

		FolderNameFromClient := subFolderDataJson["Name"].(string)

		// create the folder at this endpoint
		err := ftpClient.MakeDir(Conf.Conf.DataDirectory + "/" + nodeName + "/" + FolderNameFromClient)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			ws.Close()
			logger.WriteERRLog("api/ftp/upload.go 204 " + err.Error())
			return
		}

		// till here the client is waiting for us to initiate to send the folder's data or if the folder is empy then the client will send event "folder_empy".

		type ServerResponse struct {
			Event      string
			FolderName string
			Filename   string
		}

		serverResponse, err := json.Marshal(ServerResponse{
			Event:      "subfolder_data_start",
			FolderName: FolderNameFromClient,
			Filename:   Conf.Conf.DataFileName,
		})

		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			ws.Close()
			return
		}

		// open the file to write [Not in FTP]
		// WriteAbleFile, _ := os.OpenFile(Conf.Conf.DataDirectory+"/"+nodeName+"/"+FolderNameFromClient+"/"+"data.zip", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		// defer WriteAbleFile.Close()

		ws.WriteMessage(websocket.TextMessage, serverResponse)

		// now the client got the response to start sending data
		// it will take the data.zip and read the defined chunks and send it to server
		// now we have to wait for the client in a loop

		// since this is FTP, we will have to specifiy the set of bytes
		var CurrentByte int64 = 0
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
			if subFolderChunkDataJson["Event"] == "end_s_chunk" {
				break
			}

			if subFolderChunkDataJson["Event"].(string) != "subfolder_chunk_data" {
				logger.WriteERRLog("api/ftp/upload.go 268 ")
				logger.WriteERRLog(subFolderChunkDataJson["Event"].(string))
				ws.Close()
				return
			}

			// now we have got just a chunk of data and we to to write it
			parsedChunkFromRequest, _ := base64.StdEncoding.DecodeString(subFolderChunkDataJson["Chunk"].(string))

			// write the data to file :0
			rreader := bytes.NewReader(parsedChunkFromRequest)
			ftpClient.StorFrom(Conf.Conf.DataDirectory+"/"+nodeName+"/"+FolderNameFromClient+"/"+"data.zip", rreader, uint64(CurrentByte))
			CurrentByte = CurrentByte + int64(len(parsedChunkFromRequest))
			TotalBytes += int64(len(parsedChunkFromRequest))

			// the client was waiting for server to send back the response after writing
			serverResponse, _ := json.Marshal(ServerResponse{
				Event:      "subfolder_chunk_data_ack",
				FolderName: FolderNameFromClient,
			})

			ws.WriteMessage(websocket.TextMessage, serverResponse)
			// clean memory
			parsedChunkFromRequest = nil
			rreader = nil
			subFolderChunkData = nil
			subFolderChunkDataJson = nil
		}

	}

}

func ConnectFtp() (*ftp.ServerConn, error) {
	ftpClient, err := ftp.Dial(Conf.Conf.Ftp.FtpUrl, ftp.DialWithTimeout(5*time.Second))

	if err != nil {
		return nil, err
	}

	err = ftpClient.Login(Conf.Conf.Ftp.User, Conf.Conf.Ftp.Pass)
	if err != nil {
		return nil, err
	}

	return ftpClient, nil
}
