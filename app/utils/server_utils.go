package utils

import (
	"fmt"
	"runtime"
	"strings"
)

//Global Scope Variables
const PathToCss = "static/css"

const PathToModels = "static/models"

const PathToScripts = "static/scripts"

const PathToImages = "static/images"

const PathToSvg = "static/svg"

//returns the root path directory for the working application
func GetLocalEnv() string {
	_, filename, _, _ := runtime.Caller(0)
	tmp := strings.Split(filename, "/")
	var path string
	for i := 0; i < len(tmp)-3; i++ {
		path += tmp[i] + "/"
	}
	fmt.Println(path)
	return path
}
