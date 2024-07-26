package ui

import (
	"log"
	"os"
)

func Appendtotextfile(fn string, m string) error {
	var err error
	var path string
	path, err = os.Getwd()
	if err != nil {
		panic("path")
	}
	if path[len(path)-2:] == "ui" {
		path = path[:len(path)-3]
	}
	if path[len(path)-4:] != "data" {
		path = path + string(os.PathSeparator) + "data"
	}
	if _, err = os.Stat(path); err != nil {
		log.Println("#1 Adding folder data: " + path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err != nil {
				panic(err.Error())
			}
			// file does not exist
		} else {
			panic(err.Error())
			// other error
		}
	}
	fn = path + string(os.PathSeparator) + fn
	d1 := []byte(m)
	err = os.WriteFile(fn, d1, 0644)
	if err != nil {
		panic(err.Error())
	}
	return err
}
