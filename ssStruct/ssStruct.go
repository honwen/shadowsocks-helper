package ssStruct

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

type pair struct {
	a, b interface{}
}

var (
	testCases = []pair{
		// pair{"http://dl.google.com/pinyin/v2/GooglePinyinInstaller.exe", 60 * time.Second}, /*size: 15.7M, AtLeast 261KB/s*/
		pair{"http://go.microsoft.com/fwlink/?LinkId=691209", 60 * time.Second},         /*size: 17.4M, AtLeast 290KB/s*/
		pair{"http://dl.google.com/dl/picasa/gpautobackup_setup.exe", 20 * time.Second}, /*size: 2.5M,  AtLeast 125KB/s*/
		pair{"http://www.google.com/robots.txt", 3 * time.Second},                       /*size: 6.3K,  AtLeast 2.1KB/s*/
	}
)

// TestCase Level
type TestCase int

// TestCase Level
const (
	TestCaseMid   TestCase = 0
	TestCaseSmall TestCase = 1
	TestCaseTiny  TestCase = 2
)

// Regexp to Check URI
var (
	RegxIsSSURI  = regexp.MustCompile(`ss:\/\/([a-zA-Z0-9\.\-_]+):([^ @]+)@([a-zA-Z0-9\.\-_]+):(\d+)`)
	RegxIsSSRURI = regexp.MustCompile(`ssr:\/\/([a-zA-Z0-9\.\-_]+:)(\d+:)([a-zA-Z0-9\.\-_]+:)+\S+`)
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

// SliceSS is Slice of SS
type SliceSS []SS

// SliceSSR is Slice of SSR
type SliceSSR []SSR

// SliceFuncSSR is Slice of FuncSSR
type SliceFuncSSR []FuncSSR

func (sfssr SliceFuncSSR) Len() int {
	return len(sfssr)
}
func (sfssr SliceFuncSSR) Swap(i, j int) {
	sfssr[i], sfssr[j] = sfssr[j], sfssr[i]
}
func (sfssr SliceFuncSSR) Less(i, j int) bool {
	return sfssr[i].Speed < sfssr[j].Speed
}

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
	if err = cmd.Start(); err != nil {
		return
	}
	// time.AfterFunc(timeout+3*time.Second, func() { cmd.Process.Kill() })
	tmBegin := time.Now()
	err = errors.New("Just Focre to wGet")
	for i := 0; i < 2000 /*10s*/ && err != nil; i++ {
		time.Sleep(50 * time.Millisecond) // Wait For SSR ready
		tmdelay = time.Now().Sub(tmBegin)
		bs, err = wGetRawFastBySOCKS5Proxy(uriAddr, fmt.Sprintf("127.0.0.1:%d", randPort), timeout)
		if time.Now().Sub(tmBegin) > timeout+10*time.Second {
			break // Max Retry Time = timeout + max(socks ready delay)
		}
	}

	if err = cmd.Process.Signal(os.Kill); err != nil {
		cmd.Process.Kill()
	}
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
			Proxy:                 nil,
			Dial:                  proxy.Dial,
			ResponseHeaderTimeout: 3 * time.Second,
			IdleConnTimeout:       timeout,
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

// SpeedTest Test Speed of SSR Proxy
func (ssr *FuncSSR) SpeedTest(url string, timeout time.Duration) error {
	tmBegin := time.Now()
	if bytes, tmdelay, err := ssr.WGet(url, timeout); err == nil {
		tmTotal := time.Now().Sub(tmBegin) - tmdelay
		ssr.Speed = float64(len(bytes)) / float64(tmTotal.Nanoseconds()) * 1000.0 // Bytes / Nanoseconds * 1000 => MB/s
		if ssr.Speed < 0.0005 {
			ssr.Speed = -1
		}
	} else {
		ssr.Speed = -1
		return err
	}
	return nil
}

// SpeedTest Test Speed of SSR Proxy
func (sfssr *SliceFuncSSR) SpeedTest(TestCaseLevel TestCase) {
	var wg sync.WaitGroup
	for idx := range *sfssr {
		wg.Add(1)
		go func(idx int) {
			log.Println(fmt.Sprintf("#%02d", idx), "Test begin:", (*sfssr)[idx].SSR.String())
			(*sfssr)[idx].SpeedTest(
				testCases[TestCaseLevel].a.(string),
				testCases[TestCaseLevel].b.(time.Duration))
			log.Println(fmt.Sprintf("#%02d", idx), "Test done:", fmt.Sprintf("Speed:%6.03fMB/s", (*sfssr)[idx].Speed))
			wg.Done()
		}(idx)
		time.Sleep(333 * time.Millisecond)
		// fmt.Println(sfssr.countSpeedTested(), idx)
		for sfssr.countSpeedTested()+4 <= idx /*Max 5 Server at the same time*/ {
			// fmt.Println(sfssr.countSpeedTested(), idx)
			time.Sleep(333 * time.Millisecond)
		}
	}
	wg.Wait()
	if !sort.IsSorted(sort.Reverse(sfssr)) {
		sort.Sort(sort.Reverse(sfssr))
	}
}

func (sfssr SliceFuncSSR) countSpeedTested() (count int) {
	for _, s := range sfssr {
		if s.Speed != 0 {
			count++
		}
	}
	return
}
