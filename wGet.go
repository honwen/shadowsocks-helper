package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

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
