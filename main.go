package main

import (
	"errors"
	"fmt"
)

var (
	ErrNoGuiConfigs    = errors.New(`NO gui-config*.json`)
	ErrWrongGuiConfigs = errors.New(`Wrong gui-config*.json`)

	guiConfig *Config
)

func main() {
	if files := getFileList(getCurrDir(), isGuiConfig); files != nil {
		if guiConfig, _ = ParseConfig(files[0]); guiConfig == nil {
			panic(ErrWrongGuiConfigs)
		}
	} else {
		panic(ErrNoGuiConfigs)
	}
	fmt.Println(guiConfig)
}
