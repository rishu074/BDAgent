package api

import (
	"net/http"

	Conf "github.com/NotRoyadma/auto_backup-dnxrg/config"
	Static "github.com/NotRoyadma/auto_backup-dnxrg/routes/static"
)

func UploadFileManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("token") != Conf.Conf.Token {
		Static.ErrorRouteHandler(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if r.Header.Get("node") == "" {
		http.Error(w, "Node name must be specified", http.StatusBadRequest)
		return
	}

	if r.Header.Get("urln") == "" {
		http.Error(w, "URL Must be Specified from where to download", http.StatusBadRequest)
		return
	}

}
