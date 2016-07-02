package main

import (
	"errors"
	"fmt"
)

var (
	ErrNoGuiConfigs = errors.New(`NO gui-config*.json`)

	guiConfig *Config
	gfwList   *GFWList
)

func main() {
	var err error
	if files := getFileList(getCurrDir(), isGuiConfig); files != nil {
		if guiConfig, err = ParseConfig(files[0]); err != nil {
			panic(err)
		}
	} else {
		panic(ErrNoGuiConfigs)
	}
	fmt.Println(guiConfig)

	if gfwList, err = ParseGFWList(guiConfig.GetServerArray()); err != nil {
		panic(err)
	}
	fmt.Println(gfwList.GenDnsmasqServer("127.0.0.1:5353"))
	fmt.Println(gfwList.GenDnsmasqIPset("gfwlist"))
}
