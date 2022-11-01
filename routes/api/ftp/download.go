package ftp

import (
	"io"
	"net/http"
	"strings"

	Tools "github.com/NotRoyadma/BDAgent/avails"
	Conf "github.com/NotRoyadma/BDAgent/config"
	"github.com/NotRoyadma/BDAgent/logger"
	Static "github.com/NotRoyadma/BDAgent/routes/static"
)

func DowloadFileManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	arrayOfPath := strings.Split(r.URL.Path, "/")

	if len(arrayOfPath) != 5 {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	nodeName := arrayOfPath[3]
	servername := arrayOfPath[4]

	if !Tools.StringInSlice(nodeName, Conf.Conf.Nodes) {
		Static.ErrorRouteHandler(w, r, "No Nodes Found with this name.", 404)
		return
	}

	ftpClient, err := ConnectFtp()
	defer ftpClient.Quit()

	if err != nil {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.WriteERRLog(" /api/ftp/download.go 39 " + err.Error())
		return
	}

	if DoDirExists, _ := Tools.DoDirExistsFTP(ftpClient, "/", Conf.Conf.DataDirectory); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "No Folder found, probably due to misconfiguratin", 404)
		return
	}

	if DoDirExists, _ := Tools.DoDirExistsFTP(ftpClient, "/"+Conf.Conf.DataDirectory+"/", nodeName); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "No File found, probably due to misconfiguratin", 404)
		return
	}

	if DoDirExists, _ := Tools.DoDirExistsFTP(ftpClient, Conf.Conf.DataDirectory+"/"+nodeName, servername); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "Sorry, but your server not Found", 404)
		return
	}

	if DoDirExists, _ := Tools.DoDirExistsFTP(ftpClient, Conf.Conf.DataDirectory+"/"+nodeName+"/"+servername, Conf.Conf.DataFileName); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "Your server found, but it seems there was no data in your server.", 404)
		return
	}

	_ = ftpClient.ChangeDirToParent()

	// File found lmao
	rF, err := ftpClient.Retr("/" + Conf.Conf.DataDirectory + "/" + nodeName + "/" + servername + "/" + Conf.Conf.DataFileName)
	if err != nil {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.WriteERRLog(" /api/ftp/download.go 69 " + err.Error())
		return
	}

	//serve file
	io.Copy(w, rF)
	// http.ServeFile(w, r, Conf.Conf.DataDirectory+"/"+nodeName+"/"+servername+"/"+Conf.Conf.DataFileName)
}

// func ConnectFtp() (*ftp.ServerConn, error) {
// 	ftpClient, err := ftp.Dial(Conf.Conf.Ftp.FtpUrl, ftp.DialWithTimeout(5*time.Second))

// 	if err != nil {
// 		return nil, err
// 	}

// 	err = ftpClient.Login(Conf.Conf.Ftp.User, Conf.Conf.Ftp.Pass)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ftpClient, nil
// }
