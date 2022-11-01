package Default

import (
	"net/http"
	"strings"

	Conf "github.com/NotRoyadma/auto_backup-dnxrg/config"
	Logger "github.com/NotRoyadma/auto_backup-dnxrg/logger"
	Api "github.com/NotRoyadma/auto_backup-dnxrg/routes/api"
	Ftp "github.com/NotRoyadma/auto_backup-dnxrg/routes/api/ftp"
	Static "github.com/NotRoyadma/auto_backup-dnxrg/routes/static"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Write the HTTP logs
	Logger.WriteAutoHTTPLogs(w, r)

	//Handle Different Paths
	if r.URL.Path == "/" {
		Static.IndexRouter(w, r)
		return
	} else if r.URL.Path == "/api/status" {
		Api.StatusApiHandler(w, r)
		return
	} else if strings.Contains(r.URL.Path, "/api/download/") {
		if Conf.Conf.Ftp.Enabled {
			Ftp.DowloadFileManager(w, r)
			return
		} else {
			Api.DowloadFileManager(w, r)
			return
		}
	} else if strings.Contains(r.URL.Path, "/api/upload") {
		if Conf.Conf.Ftp.Enabled {
			Ftp.UploadFileManager(w, r)
			return
		} else {
			Api.UploadFileManager(w, r)
			return
		}
	}

	//404 page
	Static.ErrorRouteHandler(w, r, "Not Found", 404)
}
