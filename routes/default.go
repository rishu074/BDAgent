package Default

import (
	"net/http"

	Api "github.com/NotRoyadma/auto_backup-dnxrg/routes/api"
	Static "github.com/NotRoyadma/auto_backup-dnxrg/routes/static"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		Static.IndexRouter(w, r)
		return
	} else if r.URL.Path == "/api/status" {
		Api.StatusApiHandler(w, r)
		return
	}

	//404 page
	Static.ErrorRouteHandler(w, r, "Not Found", 404)
}
