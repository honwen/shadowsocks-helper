package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/chenhw2/shadowsocks-helper/ssStruct"
)

var (
	isSSRForm = regexp.MustCompile(`ssr:\/\/([a-zA-Z0-9\.\-_]+:)(\d+:)([a-zA-Z0-9\.\-_]+:)+[a-zA-Z0-9]+`)
)

// ParseSSRFromTEXT Parse ssStruct.SliceSSR From [ssr-list.txt]
func ParseSSRFromTEXT(path string) (ssrs ssStruct.SliceSSR, err error) {
	var (
		file  *os.File
		bytes []byte
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	if bytes, err = ioutil.ReadAll(file); err != nil {
		return
	}
	ssrListStr := ssStruct.RegxIsSSRURI.FindAllString(string(bytes), -1)

	for _, line := range ssrListStr {
		if strings.Count(line, `:`) >= 5 {
			line = strings.TrimPrefix(line, `ssr://`)
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")
			sub := strings.Split(line, ":")
			ssrs = append(ssrs, ssStruct.SSR{
				SS: ssStruct.SS{
					Server:     sub[0],
					ServerPort: json.Number(sub[1]),
					Method:     sub[3],
					Password:   sub[5],
				},
				Protocol: sub[2],
				Obfs:     sub[4],
			})
		} else {
			// TODO add full base64 SSR-QRcode-Scheme Support
		}
	}
	return
}

// ParseSSRFromJSON Parse ssStruct.SliceSSR From [gui-config.json]
func ParseSSRFromJSON(path, protocol, obfs string) (sss ssStruct.SliceSSR, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	mGuiConfig := &struct {
		Servers ssStruct.SliceSS `json:"configs"`
	}{}

	err = json.Unmarshal(data, mGuiConfig)
	if err != nil {
		return
	}
	sss = make(ssStruct.SliceSSR, len(mGuiConfig.Servers))
	for idx := range sss {
		sss[idx].SS = mGuiConfig.Servers[idx]
		sss[idx].Protocol = protocol
		sss[idx].Obfs = obfs
	}
	return
}

// ParseSSFromJSON Parse ssStruct.SliceSS From [gui-config.json]
func ParseSSFromJSON(path string) (sss ssStruct.SliceSS, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	mGuiConfig := &struct {
		Servers ssStruct.SliceSS `json:"configs"`
	}{}

	err = json.Unmarshal(data, mGuiConfig)
	sss = mGuiConfig.Servers
	if err != nil {
		sss = nil
	}
	return
}
