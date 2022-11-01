package logger

import (
	"log"
	"net/http"
	"os"
	"time"

	Conf "github.com/NotRoyadma/auto_backup-dnxrg/config"
)

func WriteLog(message string) bool {
	// open the log file if not create it
	LoggerFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	if err != nil {
		return false
	}

	defer LoggerFile.Close()
	n, err := LoggerFile.Write([]byte("\n" + time.Now().String() + " > " + message))
	if err != nil || n == 0 {
		return false
	}

	log.Println("[APP LOG] > ", message)

	return true
}

func WriteERRLog(message string) bool {
	// open the log file if not create it
	LoggerFile, err := os.OpenFile("./logs/app.error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	if err != nil {
		return false
	}

	defer LoggerFile.Close()
	n, err := LoggerFile.Write([]byte("\n" + time.Now().String() + " > " + message))
	if err != nil || n == 0 {
		return false
	}

	log.Println("[APP ERROR LOG] > ", message)

	return true
}

func WriteHttpLogs(message string) bool {
	// open the log file if not create it
	LoggerFile, err := os.OpenFile("./logs/http.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	if err != nil {
		return false
	}

	defer LoggerFile.Close()
	n, err := LoggerFile.Write([]byte("\n" + time.Now().String() + " > " + message))
	if err != nil || n == 0 {
		return false
	}

	log.Println("[HTTP LOG] > ", message)

	return true
}

func WriteAutoHTTPLogs(w http.ResponseWriter, r *http.Request) bool {
	// open the log file if not create it
	LoggerFile, err := os.OpenFile("./logs/http.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return false
	}

	defer LoggerFile.Close()

	var ip string
	// get the ip adress
	if Conf.Conf.IpHeader != "default" {
		ip = r.Header.Get(Conf.Conf.IpHeader)
	} else {
		ip = r.RemoteAddr
	}

	//get additional details
	path := r.URL.Path
	userAgent := r.UserAgent()
	token := r.Header.Get("token")

	//log it to file
	n, err := LoggerFile.Write([]byte("\n" + time.Now().String() + " > " + ip + " " + userAgent + " " + path + " " + token))
	if err != nil || n == 0 {
		return false
	}

	log.Println("[HTTP LOG] > ", ip+" "+userAgent+" "+path+" "+token)
	return true

}

func DeleteLogFiles() {
	er := os.RemoveAll("./logs/")
	if er != nil {
		log.Panic(er)
	}

	os.Mkdir("./logs", 0777)
}
