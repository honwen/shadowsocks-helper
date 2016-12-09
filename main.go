package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chenhw2/shadowsocks-helper/ssStruct"
	"github.com/urfave/cli"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))

	// Handle os.Signal: os.Interrupt, syscall.SIGTERM
	cleanup := make(chan os.Signal, 2)
	signal.Notify(cleanup, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cleanup
		syscall.Kill(-os.Getpid(), syscall.SIGKILL)
	}()
}

var version = "MISSING build version [git hash]"

func main() {
	app := cli.NewApp()
	app.Name = "ss(r)-helper"
	app.Usage = "shadowsocks(R)-helper"
	app.Version = version
	app.Commands = []cli.Command{
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
				// fmt.Println("ssr2json task: ", c.String("ssr"), c.String("protocol"), c.String("obfs"))
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
					Usage: "`PATH` of ssr-local (Need if proxy is not empty)",
				}},
			Action: func(c *cli.Context) error {
				// fmt.Println("dnsmasq task: ", c.String("server"), c.String("ipset"), c.String("proxy"), c.String("path"))
				var (
					base64Str string
					err       error
				)
				if len(c.String("proxy")) == 0 {
					base64Str, err = wGet(officalGFWListURL)
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
					bs, tm, err := funcSSR.WGet(officalGFWListURL, 10*time.Second)
					fmt.Println("Delay Time:", tm)
					base64Str = string(bs)
				}
				if err != nil {
					log.Printf("%+v", err)
				} else {
					gfwlist, err := ParseGFWList(base64Str)
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
					Usage: "`PATH` of ssr-local (Need if proxy is not empty)",
				}},
			Action: func(c *cli.Context) error {
				fmt.Println("dnsmasq task: ", c.String("input"), c.String("mode"), c.String("justone"), c.String("path"))
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
					bytes, err := readFileAll(c.String("justone"))
					if err != nil {
						log.Printf("%+v", err)
						return nil
					}
					ssr, err := ssStruct.ParseSSRFromURI(string(bytes))
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
					fssr.SpeedTest(ssTC)
				}
				for idx := range fssr {
					fmt.Printf("Speed:%6.03fMB/s, %s\n", fssr[idx].Speed, fssr[idx])
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
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
