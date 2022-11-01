package avails

import (
	"log"
	"strings"

	ftp "github.com/jlaffaye/ftp"
)

func DoDirExistsFTP(client *ftp.ServerConn, path string, dir string) (bool, error) {
	Directories, err := client.List(path)
	if err != nil {
		return false, err
	}

	dir = strings.TrimPrefix(dir, "./")
	log.Println(dir)

	for _, d := range Directories {
		if d.Name == dir {
			return true, nil
		}
	}

	return false, nil
}
