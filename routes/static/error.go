package static

import (
	"net/http"
)

func ErrorRouteHandler(w http.ResponseWriter, r *http.Request, msg string, code int) {
	http.Error(w, msg, code)
}
