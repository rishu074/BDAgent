package static

import (
	"net/http"
)

func IndexRouter(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}
