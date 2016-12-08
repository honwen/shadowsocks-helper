package ssStruct

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

// RegxIsSSRURI is regexp aht check ssr URI
var RegxIsSSRURI = regexp.MustCompile(`ssr:\/\/([a-zA-Z0-9\.\-_]+:)(\d+:)([a-zA-Z0-9\.\-_]+:)+[a-zA-Z0-9]+`)

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

// CmdSSR implement func
type CmdSSR struct {
	SS
	Speed float64
	cmd   string
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
		return nil, errors.New("Not SSR Form[ssr://host:port:protocol:method:obfs:pass]")
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
