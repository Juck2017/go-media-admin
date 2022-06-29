package common

import (
	"os"
)

func GetBannerString() string {
	file, err := os.Open("./banner.txt")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			GlobalConf.Logger.Error("Close Banner file error: ", err)
			panic(err)
		}
	}(file)
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		return "Bad Error"
	}
	return string(buf[:n])
}
