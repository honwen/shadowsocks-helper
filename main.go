package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	cidrman "github.com/EvilSuperstars/go-cidrman"
	"github.com/chenhw2/go-ps" /*ps*/
	"github.com/chenhw2/shadowsocks-helper/cidr"
	"github.com/chenhw2/shadowsocks-helper/ssStruct"
	"github.com/chenhw2/shadowsocks-helper/subscribe"
	"gopkg.in/urfave/cli.v1"
)

const ipipCHN = "https://raw.githubusercontent.com/17mon/china_ip_list/master/china_ip_list.txt"

var (
	version        = "MISSING build version [git hash]"
	defaultSSRPath = "ssr-local"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))

	// Handle os.Signal: os.Interrupt, syscall.SIGTERM
	cleanup := make(chan os.Signal, 2)
	signal.Notify(cleanup, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cleanup
		psClean(os.Getpid())
		os.Exit(-1)
	}()

	switch runtime.GOOS {
	case "linux":
		defaultSSRPath = "/usr/bin/ssr-local"
	case "windows":
		defaultSSRPath = getCurrDir("ssr-local.exe")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "ss(r)-helper"
	app.Usage = "shadowsocks(R)-helper"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:     "subscribe2ssr",
			Category: "CONVERTER",
			Usage:    "convert FROM[subscribe_url] to URI[ssr://host:port:protocol:method:obfs:pass]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "url, s",
					Usage: "specify `URL` of Subscribe",
				},
				cli.BoolFlag{
					Name:  "test, t",
					Usage: "SpeedTest of each ssr",
				},
			},
			Action: func(c *cli.Context) error {
				// fmt.Println("subscribe2ssr task: ", c.String("url"))
				ssrURIs, ssrRemarks := subscribe.URL2URIs(c.String("url"))
				for i := 0; i < len(ssrURIs); i++ {
					ssr := ssrURIs[i]
					remark := ssrRemarks[i]
					remark = regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(remark, "")
					remark += regexp.MustCompile(`.*@([^:]*):.*`).ReplaceAllString(ssr, `_$1`)
					fmt.Printf("%-24s : %s", remark, ssr)
					if c.Bool("t") {
						if data, t, err := subscribe.WGetRawFastBySSRProxy(`https://myip.ipip.net`, ssr, 15*time.Second); err == nil {
							fmt.Printf(` # %s : %s`, strings.ReplaceAll(string(data), "\n", ""), t.String())
						} else {
							fmt.Printf(` # BAD`)
						}
					}
					fmt.Println(``)

				}
				return nil
			},
		},
		{
			Name:     "ssr2clash",
			Category: "CONVERTER",
			Usage:    "convert FROM URI[ssr://host:port:protocol:method:obfs:pass] to CLASH config",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "url, s",
					Usage: "specify `URL` of Subscribe",
				},
			},
			Action: func(c *cli.Context) error {
				// fmt.Println("ssr2clash task: ", c.String("url"))
				subscribe.URL2Clash(c.String("url"))
				return nil
			},
		},
		{
			Name:     "json2ssr",
			Category: "CONVERTER",
			Usage:    "convert FROM[gui-config.json] to URI[ssr://host:port:protocol:method:obfs:pass]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "gui-config.json",
					Usage: "specify `PATH` of gui-config.json",
				},
				cli.StringFlag{
					Name:  "protocol, O",
					Value: "origin",
					Usage: "protocol `PLUGIN`: [origin, verify_sha1, auth_sha1_v2, auth_sha1_v4, auth_aes128_md5, auth_aes128_sha1]",
				},
				cli.StringFlag{
					Name:  "obfs, o",
					Value: "plain",
					Usage: "obfs `PLUGIN`: [plain, http_simple, http_post, tls1.2_ticket_auth]",
				},
			},
			Action: func(c *cli.Context) error {
				// fmt.Println("gui2ssr task: ", c.String("input"), c.String("protocol"), c.String("obfs"))
				ssrs, err := ParseSSRFromJSON(c.String("input"), c.String("protocol"), c.String("obfs"))
				if err != nil {
					log.Printf("%+v", err)
				}
				for _, ssr := range ssrs.GenSSRURIList() {
					fmt.Println(ssr)
				}
				return nil
			},
		},
		{
			Name:     "json2ss",
			Category: "CONVERTER",
			Usage:    "convert FROM[gui-config.json] to URI[ss://method:pass@host:port]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "gui-config.json",
					Usage: "specify `PATH` of gui-config.json",
				}},
			Action: func(c *cli.Context) error {
				// fmt.Println("gui2ss task: ", c.String("input"))
				sss, err := ParseSSFromJSON(c.String("input"))
				if err != nil {
					log.Printf("%+v", err)
				}
				for _, ss := range sss.GenSSURIList() {
					fmt.Println(ss)
				}
				return nil
			},
		},
		{
			Name:     "ss2ssr",
			Category: "CONVERTER",
			Usage:    "convert URI[ss://method:pass@host:port] to URI[ssr://host:port:protocol:method:obfs:pass]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ss, s",
					Value: "ss://method:pass@host:port",
					Usage: "shadowsocks `URI`",
				},
				cli.StringFlag{
					Name:  "protocol, O",
					Value: "origin",
					Usage: "protocol `PLUGIN`: [origin, verify_sha1, auth_sha1_v2, auth_sha1_v4, auth_aes128_md5, auth_aes128_sha1]",
				},
				cli.StringFlag{
					Name:  "obfs, o",
					Value: "plain",
					Usage: "obfs `PLUGIN`: [plain, http_simple, http_post, tls1.2_ticket_auth]",
				},
			},
			Action: func(c *cli.Context) error {
				// fmt.Println("ss2ssr task: ", c.String("ss"), c.String("protocol"), c.String("obfs"))
				ss, err := ssStruct.ParseSSFromURI(c.String("ss"))
				if err != nil {
					log.Printf("%+v", err)
				} else {
					fmt.Println(ss.ToSSR(c.String("protocol"), c.String("obfs")))
				}
				return nil
			},
		},
		{
			Name:     "ss2json",
			Category: "CONVERTER",
			Usage:    "convert URI[ss://method:pass@host:port] to JSON",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ss, s",
					Value: "ss://method:pass@host:port",
					Usage: "shadowsocks `URI`",
				},
			},
			Action: func(c *cli.Context) error {
				// fmt.Println("ss2json task: ", c.String("ss"), c.String("protocol"), c.String("obfs"))
				ss, err := ssStruct.ParseSSFromURI(c.String("ss"))
				if err != nil {
					log.Printf("%+v", err)
				} else {
					fmt.Println(ss.JSON())
				}
				return nil
			},
		},
		{
			Name:     "ssr2json",
			Category: "CONVERTER",
			Usage:    "convert URI[ssr://host:port:protocol:method:obfs:pass] to JSON",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ssr, s",
					Value: "ssr://host:port:protocol:method:obfs:pass",
					Usage: "shadowsocksR `URI`",
				},
			},
			Action: func(c *cli.Context) error {
				// fmt.Println("ssr2json task: ", c.String("ssr"))
				ssr, err := ssStruct.ParseSSRFromURI(c.String("ssr"))
				if err != nil {
					log.Printf("%+v", err)
				} else {
					fmt.Println(ssr.JSON())
				}
				return nil
			},
		},
		{
			Name:     "cidr",
			Category: "HELPER",
			Usage:    "generate CIDRs of google/netflix, etc",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "target, s",
					Value: "google",
				}},
			Action: func(c *cli.Context) error {
				cidrs := []string{}
				switch c.String("target") {
				case "google":
					cidrs, _ = cidrman.MergeCIDRs(cidr.Google())
				}
				for _, v := range cidrs {
					fmt.Println(v)
				}
				return nil
			},
		},
		{
			Name:     "dnsmasq",
			Category: "HELPER",
			Usage:    "generate DNSMASQ(server/ipset) from GFWLIST(online) with SSR-Proxy(optional)",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "server, s",
					Value: "208.67.222.220#443",
					Usage: "`DNS` for dnsmasq FROM[host#port]",
				},
				cli.StringFlag{
					Name:  "ipset, i",
					Value: "shadowsocks",
					Usage: "`NAME` for ipset, set NULL to disable",
				},
				cli.StringFlag{
					Name:  "proxy, p",
					Usage: "only support [ssr://host:port:protocol:method:obfs:pass]",
				},
				cli.StringFlag{
					Name:  "path, b",
					Value: defaultSSRPath,
					Usage: "`PATH` of ssr-local (Need if proxy is not empty)",
				}},
			Action: func(c *cli.Context) error {
				// fmt.Println("dnsmasq task: ", c.String("server"), c.String("ipset"), c.String("proxy"), c.String("path"))
				var (
					base64Str string
					extraList string
					err       error
				)
				if len(c.String("proxy")) == 0 {
					base64Str, err = wGet(officalGFWListURL)
					if err == nil {
						extraList, err = wGet(officalGoogleDomain)
					}
					extraList += strings.Replace(extraList, "google", "blogspot", -1)
				} else {
					ssr, err := ssStruct.ParseSSRFromURI(c.String("proxy"))
					if err != nil {
						log.Printf("%+v", err)
						return nil
					}
					funcSSR := ssStruct.FuncSSR{
						SSR:  *ssr,
						Path: c.String("path"),
					}
					gfwb64, tm, err := funcSSR.WGet(officalGFWListURL, 10*time.Second)
					fmt.Println("Delay Time:", tm)
					base64Str = string(gfwb64)
					if err == nil {
						bs, _, _ := funcSSR.WGet(officalGoogleDomain, 10*time.Second)
						extraList = string(bs)
						extraList += strings.Replace(extraList, "google", "blogspot", -1)
						fmt.Println(extraList)
					}
				}
				if err != nil {
					log.Printf("%+v", err)
				} else {
					gfwlist, err := ParseGFWList(base64Str, extraList)
					if err != nil {
						log.Printf("%+v", err)
					} else {
						fmt.Println(gfwlist.GenDnsmasqServer(c.String("server")))
						if !strings.EqualFold(c.String("ipset"), "NULL") {
							fmt.Println(gfwlist.GenDnsmasqIPset(c.String("ipset")))
						}
					}
				}
				return nil
			},
		},
		{
			Name:     "ssrrank",
			Category: "HELPER",
			Usage:    "speed-test for shadowsocksr from list of FROM[ssr://host:port:protocol:method:obfs:pass]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "ssr-list.txt",
					Usage: "specify `PATH` of ssr-list.txt",
				},
				cli.StringFlag{
					Name:  "mode, m",
					Value: "tiny",
					Usage: "`Mode` for speed-test [mid, small, tiny]",
				},
				cli.StringFlag{
					Name:  "justone",
					Usage: "only test specify the ssr://host:port:protocol:method:obfs:pass, this will disable input",
				},
				cli.StringFlag{
					Name:  "path, b",
					Value: defaultSSRPath,
					Usage: "`PATH` of ssr-local",
				}},
			Action: func(c *cli.Context) error {
				// fmt.Println("ssrrank task: ", c.String("input"), c.String("mode"), c.String("justone"), c.String("path"))
				var ssTC ssStruct.TestCase
				var fssr ssStruct.SliceFuncSSR
				switch c.String("mode") {
				case "mid":
					ssTC = ssStruct.TestCaseMid
				case "small":
					ssTC = ssStruct.TestCaseSmall
				case "tiny":
					fallthrough
				default:
					ssTC = ssStruct.TestCaseTiny
				}
				if len(c.String("justone")) > 0 {
					fmt.Println(c.String("justone"))
					ssr, err := ssStruct.ParseSSRFromURI(c.String("justone"))
					if err != nil {
						log.Printf("%+v", err)
						return nil
					}
					fssr = ssStruct.SliceFuncSSR{
						ssStruct.FuncSSR{
							SSR:   *ssr,
							Path:  c.String("path"),
							Speed: 0,
						},
					}
					fssr.SpeedTest(ssTC)
				} else {
					ssr, err := ParseSSRFromTEXT(c.String("input"))
					if err != nil {
						log.Printf("%+v", err)
						return nil
					}
					fssr = make(ssStruct.SliceFuncSSR, len(ssr))
					for idx := range ssr {
						fssr[idx] = ssStruct.FuncSSR{
							SSR:  ssr[idx],
							Path: c.String("path"),
						}
					}
					// fmt.Println("Total SSR Server:", fssr.Len(), "\n", fssr)
					fssr.SpeedTest(ssTC)
				}
				for idx := range fssr {
					fmt.Printf("Speed:%6.03fMB/s, %s\n", fssr[idx].Speed, fssr[idx])
				}
				return nil
			},
		},
		{
			Name:     "ipset",
			Category: "HELPER",
			Usage:    "Minify chnroute.txt",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "chnroute.txt",
					Usage: "specify `PATH` of chnroute.txt",
				},
				cli.BoolFlag{
					Name:  "online, n",
					Usage: "Get chnroute.txt from [IPIP.NET]",
				},
				cli.StringFlag{
					Name:  "proxy, p",
					Usage: "only support [ssr://host:port:protocol:method:obfs:pass]",
				},
				cli.Int64Flag{
					Name:  "maxMask, m",
					Value: 16,
					Usage: "Like: 192.168.1.1/13",
				}},
			Action: func(c *cli.Context) error {
				var (
					chnroute string
					err      error
				)
				if c.BoolT("online") {
					if len(c.String("proxy")) == 0 {
						chnroute, err = wGet(ipipCHN)
					} else {
						ssr, err := ssStruct.ParseSSRFromURI(c.String("proxy"))
						if err != nil {
							log.Printf("%+v", err)
							return nil
						}
						funcSSR := ssStruct.FuncSSR{
							SSR:  *ssr,
							Path: c.String("path"),
						}
						bytes, tm, err := funcSSR.WGet(ipipCHN, 10*time.Second)
						fmt.Println("Delay Time:", tm)
						chnroute = string(bytes)
					}
				} else {
					var bytes []byte
					if bytes, err = readFileAll(c.String("input")); err != nil {
						return err
					}
					chnroute = string(bytes)
				}
				maxMask := c.Int64("maxMask")
				if maxMask > 31 {
					maxMask = 31
				}
				if maxMask < 5 {
					maxMask = 5
				}
				// fmt.Println(chnroute)
				for i := int64(32); i > maxMask; i-- {
					chnroute = regexp.MustCompile(fmt.Sprintf("/%d", i)).ReplaceAllString(chnroute, fmt.Sprintf("/%d", i-1))
				}
				chnipset, err := cidrman.MergeCIDRs(strings.Split(chnroute, "\n"))
				if err != nil {
					log.Printf("%+v", err)
					return nil
				}
				for _, v := range chnipset {
					fmt.Println(v)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
	psClean(os.Getpid())
}

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

func getCurrDir(fn string) (path string) {
	path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(path, fn)
}

func psClean(pid int) {
	childPs, _ := ps.FilterProcesses(func(p ps.Process) bool {
		return p.PPid() == pid
	})
	if childPs == nil {
		if p, err := os.FindProcess(pid); err == nil {
			if err = p.Signal(os.Kill); err != nil {
				p.Kill()
			}
		}
	} else {
		for _, ps := range childPs {
			psClean(ps.Pid())
		}
	}
}
