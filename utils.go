package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

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

func ss2json(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil || u.User == nil {
		return ``, err
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return ``, err
	}
	ssJSON := struct {
		Server     string      `json:"server"`
		ServerPort json.Number `json:"server_port"`
		Password   string      `json:"password"`
		Method     string      `json:"method"`
	}{
		Server:     host,
		ServerPort: json.Number(port),
		Method:     u.User.Username(),
	}
	ssJSON.Password, _ = u.User.Password()
	b, err := json.MarshalIndent(ssJSON, "", "  ")
	if err != nil {
		return ``, err
	}
	return string(b), nil
}

func ssr2json(s string) (string, error) {
	if !isSSRForm.MatchString(s) {
		return ``, errors.New("Not SSR Form")
	}
	s = strings.TrimPrefix(s, `ssr://`)
	sub := strings.Split(s, ":")
	ssrJSON := SSRConfig{
		Server:     sub[0],
		ServerPort: json.Number(sub[1]),
		Protocol:   sub[2],
		Method:     sub[3],
		Obfs:       sub[4],
		Password:   sub[5],
	}
	b, err := json.MarshalIndent(ssrJSON, "", "  ")
	if err != nil {
		return ``, err
	}
	return string(b), nil
}

func ss2ssr(s, protocol, obfs string) (string, error) {
	u, err := url.Parse(s)
	if err != nil || u.User == nil {
		return ``, err
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return ``, err
	}
	ssr := struct {
		Server     string      `json:"server"`
		ServerPort json.Number `json:"server_port"`
		Password   string      `json:"password"`
		Method     string      `json:"method"`
		Protocol   string      `json:"protocol"`
		Obfs       string      `json:"obfs"`
	}{
		Server:     host,
		ServerPort: json.Number(port),
		Method:     u.User.Username(),
		Protocol:   protocol,
		Obfs:       obfs,
	}
	ssr.Password, _ = u.User.Password()
	return fmt.Sprintf("ssr://%s:%s:%s:%s:%s:%s",
		ssr.Server, string(ssr.ServerPort), ssr.Protocol, ssr.Method, ssr.Obfs, ssr.Password), nil
}
