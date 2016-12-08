package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"log"

	"github.com/chenhw2/shadowsocks-helper/ssStruct"
	"github.com/urfave/cli"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))

}

func main() {
	app := cli.NewApp()
	app.Name = "ss(r)-helper"
	app.Usage = "shadowsocks(R)-helper"
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
				ssr, err := ssStruct.ParseSSRFromURI(c.String("ss"))
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
			Usage:    "Gen DNSMASQ(server/ipset) from GFWLIST(online) with SSR-Proxy(optional)",
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
				}},
			Action: func(c *cli.Context) error {
				fmt.Println("dnsmasq task: ", c.String("server"), c.String("ipset"), c.String("proxy"))
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
					Value: "large, mid, little",
					Usage: "`Mode` for speed-test",
				},
				cli.StringFlag{
					Name:  "justone",
					Usage: "only test specify the ssr://host:port:protocol:method:obfs:pass, this will disable input",
				}},
			Action: func(c *cli.Context) error {
				fmt.Println("dnsmasq task: ", c.String("server"), c.String("ipset"), c.String("proxy"))
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dryrun",
			Usage: "Skip RANK",
		},
	}
	app.Run(os.Args)

	var (
	//err error
	//	wg  sync.WaitGroup
	)

	/*
		go func() {
			defer wg.Done()
			addr := ":1088"
			mux := http.NewServeMux()

			//===================================================================DNSMASQ
			log.Println("Listen " + addr + "/dnsmasq")
			mux.HandleFunc("/dnsmasq", func(w http.ResponseWriter, r *http.Request) {
				if gfwList == nil {
					io.WriteString(w, "// NOT READY")
					io.WriteString(w, "\n")
					log.Println("Listen " + addr + "/dnsmasq -- > NOT READY")
					return
				}
				if r.ParseForm(); r.Form == nil {
					return
				}
				if r.Form["server"] != nil {
					log.Println("Listen " + addr + "/dnsmasq -- > server")
					io.WriteString(w, gfwList.GenDnsmasqServer(r.Form["server"][0]))
				}
				if r.Form["ipset"] != nil {
					log.Println("Listen " + addr + "/dnsmasq -- > ipset")
					io.WriteString(w, gfwList.GenDnsmasqIPset(r.Form["ipset"][0]))
				}
			})
			//===================================================================DNSMASQ

			if !onlySSR {
				//===================================================================SS2JSON
				log.Println("Listen " + addr + "/ss2json")
				mux.HandleFunc("/ss2json", func(w http.ResponseWriter, r *http.Request) {
					if strings.HasPrefix(r.RequestURI, "/ss2json?") {
						ss := strings.TrimPrefix(r.RequestURI, "/ss2json?")
						log.Println("Listen "+addr+"/ss2json", ss)
						if json, err := ss2json(ss); err == nil {
							io.WriteString(w, json)
							io.WriteString(w, "\n")
						} else {
							log.Println("Listen "+addr+"/ss2json, ERR:", err)
							return
						}
					}
				})
				//===================================================================SS2JSON
				//===================================================================SS-RANK
				log.Println("Listen " + addr + "/ssrank")
				mux.HandleFunc("/ssrank", func(w http.ResponseWriter, r *http.Request) {
					if guiCfgSS.LastTestTime.Equal(NoLastTest) {
						io.WriteString(w, "// NOT READY")
						io.WriteString(w, "\n")
						log.Println("Listen " + addr + "/ssrank -- > NOT READY")
						return
					}
					if strings.HasPrefix(r.RequestURI, "/ssrank?") {
						switch r.RequestURI {
						case "/ssrank?renew": // renew crontrl, Duration > 1h
							if time.Now().Sub(guiCfgSS.LastTestTime) > 1*time.Hour {
								log.Println("Listen "+addr+"/ssrank", "accept renew")
								io.WriteString(w, fmt.Sprint("#", guiCfgSS.LastTestTime, ", ACCEPT", "\n"))
								go guiCfgSS.TestServersBySS()
								guiCfgSS.LastTestTime = time.Now()
							} else {
								log.Println("Listen "+addr+"/ssrank", "reject renew")
								io.WriteString(w, fmt.Sprint("#", guiCfgSS.LastTestTime, ", REJECT", "\n"))
							}
						}
					} else {
						log.Println("Listen " + addr + "/ssrank")
						io.WriteString(w, fmt.Sprint("#", guiCfgSS.LastTestTime, "\n"))
						io.WriteString(w, guiCfgSS.GenTestedConfig().ssString())
					}
				})
				//===================================================================SS-RANK
			}

			//==================================================================SSR2JSON
			log.Println("Listen " + addr + "/ssr2json")
			mux.HandleFunc("/ssr2json", func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.RequestURI, "/ssr2json?") {
					ss := strings.TrimPrefix(r.RequestURI, "/ssr2json?")
					log.Println("Listen "+addr+"/ssr2json", ss)
					if json, err := ssr2json(ss); err == nil {
						io.WriteString(w, json)
						io.WriteString(w, "\n")
					} else {
						log.Println("Listen "+addr+"/ssr2json, ERR:", err)
						return
					}
				}
			})
			//==================================================================SSR2JSON

			//==================================================================SSR-RANK
			log.Println("Listen " + addr + "/ssrrank")
			mux.HandleFunc("/ssrrank", func(w http.ResponseWriter, r *http.Request) {
				if guiCfgSSR.LastTestTime.Equal(NoLastTest) {
					io.WriteString(w, "// NOT READY")
					io.WriteString(w, "\n")
					log.Println("Listen " + addr + "/ssrrank -- > NOT READY")
					return
				}
				if strings.HasPrefix(r.RequestURI, "/ssrrank?") {
					switch r.RequestURI {
					case "/ssrrank?renew": // renew crontrl, Duration > 1h
						if time.Now().Sub(guiCfgSSR.LastTestTime) > 1*time.Hour {
							log.Println("Listen "+addr+"/ssrrank", "accept renew")
							io.WriteString(w, fmt.Sprint("#", guiCfgSSR.LastTestTime, ", ACCEPT", "\n"))
							go guiCfgSSR.TestServersBySSR()
							guiCfgSSR.LastTestTime = time.Now()
						} else {
							log.Println("Listen "+addr+"/ssrrank", "reject renew")
							io.WriteString(w, fmt.Sprint("#", guiCfgSSR.LastTestTime, ", REJECT", "\n"))
						}
					}
				} else {
					log.Println("Listen " + addr + "/ssrrank")
					io.WriteString(w, fmt.Sprint("#", guiCfgSSR.LastTestTime, "\n"))
					io.WriteString(w, guiCfgSSR.GenTestedConfig().String())
				}
			})
			//==================================================================SSR-RANK
			http.ListenAndServe(addr, mux)
		}()

		wg.Add(1)
		wg.Wait()
	*/
}
