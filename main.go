package main

import (
	"errors"
	"fmt"
)

var (
	ErrNoGuiConfigs = errors.New(`NO gui-config*.json`)

	guiConfig *Config
)

func main() {
	if files := getFileList(getCurrDir(), isGuiConfig); files != nil {
		guiConfig, _ = ParseConfig(files[0])
	} else {
		panic(ErrNoGuiConfigs)
	}
	fmt.Println(guiConfig)
}
