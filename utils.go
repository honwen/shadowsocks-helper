package main

import (
	"errors"
	"net"
	"strconv"
)

func isPortAvailable(port int) bool {
	if port < 1 || port > 65535 {
		return false
	}
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	ln.Close()
	return err == nil
}

func wGetWithSSRFailsafe(urlAddr string, ssProxy []string) (text string, err error) {
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
