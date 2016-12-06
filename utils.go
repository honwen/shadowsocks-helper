package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func getCurrDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func getFileList(path string, checkFunc func(string) bool) (files []string) {
	filepath.Walk(path, func(path string, fInfo os.FileInfo, err error) error {
		if fInfo == nil {
			return err
		}
		if fInfo.IsDir() {
			return nil
		}
		if checkFunc == nil || checkFunc(path) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func isGuiConfig(s string) bool {
	return regexp.MustCompile(`gui-config.*\.json$`).MatchString(s)
}

func isSSRList(s string) bool {
	return regexp.MustCompile(`ssr-list.*\.txt$`).MatchString(s)
}

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func isPortAvailable(port int) bool {
	if port < 1 || port > 65535 {
		return false
	}
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	ln.Close()
	return err == nil
}

func wGetWithSSFailsafe(urlAddr string, ssProxy []string) (text string, err error) {
	for i := 0; i < 5 && len(text) == 0; i++ {
		text, err = wGet(urlAddr)
	}
	if len(text) == 0 && len(ssProxy) > 0 {
		for _, ss := range ssProxy {
			text, err = wGetByShadowsocksProxy(urlAddr, ss)
			if len(text) > 0 {
				err = errors.New("wGet using " + ss)
				break
			}
		}
	}
	return text, err
}
