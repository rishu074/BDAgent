package avails

import (
	"os"
	"strconv"
)

func SzCalc(path string) int64 {
	d, _ := os.ReadDir(path)
	var size int64

	for _, file := range d {
		if file.IsDir() {
			size += SzCalc(path + "/" + file.Name())
		} else {
			info, _ := file.Info()
			size = size + info.Size()
		}
	}
	return size
}

func DirSize(path string) string {
	size := SzCalc(path)

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

func Ptettier(size int64) string {
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
