package static

import (
	"html/template"
	"net/http"

	Conf "github.com/NotRoyadma/BDAgent/config"
)

func IndexRouter(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./public/index.html"))

	type DataStruct struct {
		Nodes []string
	}

	data := DataStruct{
		Nodes: Conf.Conf.Nodes,
	}
	tmpl.Execute(w, data)
}
