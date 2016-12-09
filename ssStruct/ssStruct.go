package ssStruct

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// Regexp to Check URI
var (
	RegxIsSSURI  = regexp.MustCompile(`ss:\/\/([a-zA-Z0-9\.\-_]+):([a-zA-Z0-9\.\-_]+)@([a-zA-Z0-9\.\-_]+):(\d+)`)
	RegxIsSSRURI = regexp.MustCompile(`ssr:\/\/([a-zA-Z0-9\.\-_]+:)(\d+:)([a-zA-Z0-9\.\-_]+:)+[a-zA-Z0-9]+`)
)

// SS original
type SS struct {
	Server     string      `json:"server"`
	ServerPort json.Number `json:"server_port"`
	Password   string      `json:"password"`
	Method     string      `json:"method"`
}

// SSR implement
type SSR struct {
	SS
	Protocol string `json:"protocol"`
	Obfs     string `json:"obfs"`
}

// FuncSSR implement func
type FuncSSR struct {
	SSR
	Speed float64
	Path  string
}

// SliceSSR is Slice of SSR
type SliceSSR []SSR

// SliceSS is Slice of SS
type SliceSS []SS

// JSON gen json form of SS
func (ss SS) JSON() string {
	bs, _ := json.MarshalIndent(ss, "", "  ")
	return string(bs)
}

// JSON gen json form of SSR
func (ssr SSR) JSON() string {
	bs, _ := json.MarshalIndent(ssr, "", "  ")
	return string(bs)
}

// String gen URI[ss://method:pass@host:port] of SS
func (ss SS) String() string {
	return fmt.Sprintf("ss://%s:%s@%s:%s",
		ss.Method, ss.Password, ss.Server, string(ss.ServerPort))
}

// String gen URI[ssr://host:port:protocol:method:obfs:pass] of SSR
func (ssr SSR) String() string {
	return fmt.Sprintf("ssr://%s:%s:%s:%s:%s:%s",
		ssr.Server, string(ssr.ServerPort), ssr.Protocol, ssr.Method, ssr.Obfs, ssr.Password)
}

// ToSSR gen SSR
func (ss SS) ToSSR(protocol, obfs string) SSR {
	return SSR{
		SS:       ss,
		Protocol: protocol,
		Obfs:     obfs,
	}
}

//ParseSSFromURI Parse SS From URI[ss://method:pass@host:port]
func ParseSSFromURI(str string) (*SS, error) {
	if !RegxIsSSURI.MatchString(str) {
		return nil, errors.New("Not SS URI[ss://method:pass@host:port]")
	}
	u, err := url.Parse(str)
	if err != nil || u.User == nil {
		return nil, err
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}
	ss := &SS{
		Server:     host,
		ServerPort: json.Number(port),
		Method:     u.User.Username(),
	}
	ss.Password, _ = u.User.Password()
	return ss, nil
}

//ParseSSRFromURI Parse SSR From URI [ssr://host:port:protocol:method:obfs:pass]
func ParseSSRFromURI(str string) (*SSR, error) {
	if !RegxIsSSRURI.MatchString(str) {
		return nil, errors.New("Not SSR URI[ssr://host:port:protocol:method:obfs:pass]")
	}
	sub := strings.Split(strings.TrimPrefix(str, `ssr://`), ":")
	ssr := &SSR{
		SS: SS{
			Server:     sub[0],
			ServerPort: json.Number(sub[1]),
			Method:     sub[3],
			Password:   sub[5],
		},
		Protocol: sub[2],
		Obfs:     sub[4],
	}
	return ssr, nil
}

// GenSSURIList gen SSR to URI [ss://method:pass@host:port]
func (sss SliceSS) GenSSURIList() (ssStr []string) {
	for _, ss := range sss {
		ssStr = append(ssStr, fmt.Sprintf("ss://%s:%s@%s:%s",
			ss.Method, ss.Password, ss.Server, string(ss.ServerPort)))
	}
	return
}

// GenSSRURIList gen SSR to URI [ssr://host:port:protocol:method:obfs:pass]
func (ssrs SliceSSR) GenSSRURIList() (ssrStr []string) {
	for _, ssr := range ssrs {
		ssrStr = append(ssrStr, fmt.Sprintf("ssr://%s:%s:%s:%s:%s:%s",
			ssr.Server, string(ssr.ServerPort), ssr.Protocol, ssr.Method, ssr.Obfs, ssr.Password))
	}
	return
}

// WGet By ShadowsocksR Proxy
func (ssr FuncSSR) WGet(uriAddr string, timeout time.Duration) (bs []byte, tmdelay time.Duration, err error) {
	randPort := randomTCPPort(32)
	cmd := exec.Command(ssr.Path, "-s", ssr.Server, "-p", string(ssr.ServerPort), "-m", ssr.Method, "-k", ssr.Password,
		"-O", ssr.Protocol, "-o", ssr.Obfs, "-b", "127.0.0.1", "-l", strconv.Itoa(randPort),
	)
	// fmt.Println(cmd.Args)
	if err = cmd.Start(); err != nil {
		return
	}
	tmBegin := time.Now()
	err = errors.New("Just Focre to wGet")
	for i := 0; i < 2000 /*10s*/ && err != nil; i++ {
		// fmt.Println("SSR not ready")
		time.Sleep(50 * time.Millisecond) // Wait For SSR ready
		tmdelay = time.Now().Sub(tmBegin)
		bs, err = wGetRawFastBySOCKS5Proxy(uriAddr, fmt.Sprintf("127.0.0.1:%d", randPort), timeout)
	}
	cmd.Process.Kill()
	return
}

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

func randomTCPPort(maxRetryTimes int) (randPort int) {
	for i := 0; i < maxRetryTimes; i++ {
		randPort = 10000 + rand.Int()%5*10000 + rand.Int()%3333 + 3333 // [13333, 56665]
		if ln, err := net.Listen("tcp", ":"+strconv.Itoa(randPort)); err == nil {
			ln.Close()
			break
		}
	}
	return
}
