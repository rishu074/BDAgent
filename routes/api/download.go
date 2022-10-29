package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	Tools "github.com/NotRoyadma/auto_backup-dnxrg/avails"
	Conf "github.com/NotRoyadma/auto_backup-dnxrg/config"
	Static "github.com/NotRoyadma/auto_backup-dnxrg/routes/static"
)

func DowloadFileManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	nodeName := strings.Split(r.URL.Path, "/")[3]

	if !Tools.StringInSlice(nodeName, *&Conf.Conf.Nodes) {
		Static.ErrorRouteHandler(w, r, "No Nodes Found with this name.", 404)
		return
	}

	if DoDirExists, _ := Tools.DoDirectoryExists(*&Conf.Conf.DataDirectory + "/" + nodeName); !DoDirExists {
		Static.ErrorRouteHandler(w, r, "No File found, probably due to misconfiguratin", 404)
		return
	}

	type dataStructure struct {
		Data string
	}

	fmt.Println(Conf.Conf.Nodes)

	data, _ := json.Marshal(&dataStructure{Data: "Hello from json."})

	w.Write(data)
}
