package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"

	"github.com/chenhw2/shadowsocks-helper/ssStruct"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

// Get return text of url
func wGet(urlAddr string) (string, error) {
	request, err := http.NewRequest("GET", urlAddr, nil)
	if err != nil {
		return ``, err
	}
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return ``, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), err
}

// GetByProxy return text of url visited though proxy
func wGetByHTTPProxy(urlAddr, proxyAddr string) (string, error) {
	request, _ := http.NewRequest("GET", urlAddr, nil)
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return ``, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
			ResponseHeaderTimeout: 5 * time.Second,
			DisableKeepAlives:     true,
		},
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return ``, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), err
}

// GetByProxy return text of url visited though proxy
func wGetRawFastBySOCKS5Proxy(urlAddr, proxyAddr string, timeout time.Duration) ([]byte, error) {
	request, _ := http.NewRequest("GET", urlAddr, nil)
	proxy, err := proxy.SOCKS5("tcp", proxyAddr,
		&proxy.Auth{User: "username", Password: "password"},
		&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		},
	)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			Dial:  proxy.Dial,
			ResponseHeaderTimeout: 3 * time.Second,
			DisableKeepAlives:     true,
		},
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body, err
}

// GetByShadowsocksProxy return text of URI through Shadowsocks Proxy
func wGetByShadowsocksProxy(uriAddr, ssProxyAddr string) (string, error) {
	request, err := http.NewRequest("GET", uriAddr, nil)
	if err != nil {
		return ``, err
	}
	parsedURL, err := url.Parse(uriAddr)
	if err != nil {
		return ``, err
	}
	if _, _, err := net.SplitHostPort(parsedURL.Host); err != nil {
		switch parsedURL.Scheme {
		case `https`:
			parsedURL.Host = net.JoinHostPort(parsedURL.Host, `443`)
		case `http`:
			fallthrough
		default:
			// default http
			parsedURL.Host = net.JoinHostPort(parsedURL.Host, `80`)
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(_, _ string) (net.Conn, error) {
				ssURL, err := url.Parse(ssProxyAddr)
				if err != nil || strings.ToLower(ssURL.Scheme) != `ss` {
					return nil, err
				}
				method := strings.ToLower(ssURL.User.Username())
				password, _ := ssURL.User.Password()
				cipher, err := shadowsocks.NewCipher(method, password)
				if err != nil {
					return nil, err
				}
				rawAddr, err := shadowsocks.RawAddr(parsedURL.Host)
				if err != nil {
					return nil, err
				}
				return shadowsocks.DialWithRawAddr(
					rawAddr, strings.ToLower(ssURL.Host), cipher.Copy())
			},
			ResponseHeaderTimeout: 20 * time.Second,
			DisableKeepAlives:     true,
		},
		Timeout: 30 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		return ``, err
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return ``, err
	}
	return string(body), nil
}

// GetByShadowsocksProxy return Raw Bytes of URI through Shadowsocks Proxy
func wGetRawFastByShadowsocksProxy(uriAddr, ssProxyAddr string, timeout time.Duration) ([]byte, error) {
	request, err := http.NewRequest("GET", uriAddr, nil)
	if err != nil {
		return nil, err
	}
	parsedURL, err := url.Parse(uriAddr)
	if err != nil {
		return nil, err
	}
	if _, _, err := net.SplitHostPort(parsedURL.Host); err != nil {
		switch parsedURL.Scheme {
		case `https`:
			parsedURL.Host = net.JoinHostPort(parsedURL.Host, `443`)
		case `http`:
			fallthrough
		default:
			// default http
			parsedURL.Host = net.JoinHostPort(parsedURL.Host, `80`)
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(_, _ string) (net.Conn, error) {
				ssURL, err := url.Parse(ssProxyAddr)
				if err != nil || strings.ToLower(ssURL.Scheme) != `ss` {
					return nil, err
				}
				method := strings.ToLower(ssURL.User.Username())
				password, _ := ssURL.User.Password()
				cipher, err := shadowsocks.NewCipher(method, password)
				if err != nil {
					return nil, err
				}
				rawAddr, err := shadowsocks.RawAddr(parsedURL.Host)
				if err != nil {
					return nil, err
				}
				return shadowsocks.DialWithRawAddr(
					rawAddr, strings.ToLower(ssURL.Host), cipher.Copy())
			},
			ResponseHeaderTimeout: 3 * time.Second,
			DisableKeepAlives:     true,
		},
		Timeout: timeout,
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetByShadowsocksRProxy return Raw Bytes of URI through ShadowsocksR Proxy
func wGetRawFastByShadowsocksRProxy(ssrPath, uriAddr string, ssr ssStruct.SSR, timeout time.Duration) (b []byte, err error) {
	randPort := -1
	for !isPortAvailable(randPort) {
		randPort = 10000 + rand.Int()%3*10000 + rand.Int()%33 + 333
	}
	cmd := exec.Command(ssrPath, "-s", ssr.Server, "-p", string(ssr.ServerPort), "-m", ssr.Method, "-k", ssr.Password,
		"-O", ssr.Protocol, "-o", ssr.Obfs, "-b", "127.0.0.1", "-l", strconv.Itoa(randPort),
	)
	// fmt.Println(cmd.Args)
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second) // Wait For SSR ready
	b, err = wGetRawFastBySOCKS5Proxy(uriAddr, fmt.Sprintf("127.0.0.1:%d", randPort), timeout)
	cmd.Process.Kill()
	return
}

// GetByShadowsocksRProxy return Raw Bytes of URI through ShadowsocksR Proxy
func wGetByShadowsocksRProxy(ssrPath, uriAddr string, ssr ssStruct.SSR, timeout time.Duration) (b []byte, err error) {
	randPort := -1
	for !isPortAvailable(randPort) {
		randPort = 10000 + rand.Int()%3*10000 + rand.Int()%33 + 333
	}
	cmd := exec.Command(ssrPath, "-s", ssr.Server, "-p", string(ssr.ServerPort), "-m", ssr.Method, "-k", ssr.Password,
		"-O", ssr.Protocol, "-o", ssr.Obfs, "-b", "127.0.0.1", "-l", strconv.Itoa(randPort),
	)
	// fmt.Println(cmd.Args)
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	time.Sleep(timeout) // Wait For SSR ready
	b, err = wGetRawFastBySOCKS5Proxy(uriAddr, fmt.Sprintf("127.0.0.1:%d", randPort), timeout)
	for i := 0; i < 3 && err != nil; i++ {
		time.Sleep(2 * timeout) // Wait For SSR ready
		b, err = wGetRawFastBySOCKS5Proxy(uriAddr, fmt.Sprintf("127.0.0.1:%d", randPort), timeout)
	}
	cmd.Process.Kill()
	return
}
