package main

import (
	"log"
	"net/http"
	"strconv"

	Conf "github.com/NotRoyadma/auto_backup-dnxrg/config"
	Home "github.com/NotRoyadma/auto_backup-dnxrg/routes"
)

func main() {
	http.HandleFunc("/", Home.DefaultHandler)
	log.Println("Listening on " + strconv.Itoa(*&Conf.Conf.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*&Conf.Conf.Port), nil))
}
