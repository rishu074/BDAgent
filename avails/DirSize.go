package avails

import (
	"os"
	"strconv"
)

func szCalc(path string) int64 {
	d, _ := os.ReadDir(path)
	var size int64

	for _, file := range d {
		if file.IsDir() {
			size += szCalc(path + "/" + file.Name())
		} else {
			info, _ := file.Info()
			size = size + info.Size()
		}
	}
	return size
}

func DirSize(path string) string {
	size := szCalc(path)

	switch {
	case size > 1024*1024*1024:
		return strconv.Itoa(int((size)/(1024*1024*1024))) + "GB"
	case size > 1024*1024:
		return strconv.Itoa(int(float64(size)/(1024*1024))) + "MB"
	case size > 1024:
		return strconv.Itoa(int(float64(size)/1024)) + "KB"
	default:
		return strconv.Itoa(int(size)) + "Bytes"
	}
}
