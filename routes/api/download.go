package api

import (
	"net/http"
	"strings"

	Tools "github.com/NotRoyadma/auto_backup-dnxrg/avails"
	Conf "github.com/NotRoyadma/auto_backup-dnxrg/config"
	Static "github.com/NotRoyadma/auto_backup-dnxrg/routes/static"
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

	if DoDirExists, _ := Tools.DoDirectoryExists(Conf.Conf.DataDirectory); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "No Folder found, probably due to misconfiguratin", 404)
		return
	}

	if DoDirExists, _ := Tools.DoDirectoryExists(Conf.Conf.DataDirectory + "/" + nodeName); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "No File found, probably due to misconfiguratin", 404)
		return
	}

	if DoDirExists, _ := Tools.DoDirectoryExists(Conf.Conf.DataDirectory + "/" + nodeName + "/" + servername); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "Sorry, but your server not Found", 404)
		return
	}

	if DoDirExists, _ := Tools.DoDirectoryExists(Conf.Conf.DataDirectory + "/" + nodeName + "/" + servername + "/" + Conf.Conf.DataFileName); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "Your server found, but it seems there was no data in your server.", 404)
		return
	}

	// File found lmao
	http.ServeFile(w, r, Conf.Conf.DataDirectory+"/"+nodeName+"/"+servername+"/"+Conf.Conf.DataFileName)
}
