package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli"
)

var (
	onlySSR bool
	SSRPATH string

	guiCfgSSR *Config
	guiCfgSS  Config
	gfwList   *GFWList
)

func init() {
	app := cli.NewApp()
	app.Name = "ss(r)-helper"
	app.Usage = "shadowsocks(R)-helper"
	app.Commands = []cli.Command{
		{
			Name:     "gui2ssr",
			Category: "Converter",
			Usage:    "convert FROM[gui-config.json] to FROM[ssr://host:port:protocol:method:obfs:pass]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "gui-config.json",
					Usage: "specify `PATH` of gui-config.json",
				},
				cli.StringFlag{
					Name:  "protocol, O",
					Value: "origin",
					Usage: "protocol `PLUGIN`: [origin,verify_sha1,auth_sha1_v2,auth_sha1_v4,auth_aes128_md5,auth_aes128_sha1]",
				},
				cli.StringFlag{
					Name:  "obfs, o",
					Value: "plain",
					Usage: "obfs `PLUGIN`: [plain,http_simple,http_post,tls1.2_ticket_auth]",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("gui2ssr task: ", c.String("input"), c.String("protocol"), c.String("obfs"))
				return nil
			},
		},
		{
			Name:     "ss2ssr",
			Category: "Converter",
			Usage:    "convert FROM[ss://method[-auth]:pass@host:port] to FROM[ssr://host:port:protocol:method:obfs:pass]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ss, s",
					Value: "ss://bf-cfb:test@192.168.100.1:8888",
					Usage: "shadowsocks `URI`",
				},
				cli.StringFlag{
					Name:  "protocol, O",
					Value: "origin",
					Usage: "protocol `PLUGIN`: [origin,verify_sha1,auth_sha1_v2,auth_sha1_v4,auth_aes128_md5,auth_aes128_sha1]",
				},
				cli.StringFlag{
					Name:  "obfs, o",
					Value: "plain",
					Usage: "obfs `PLUGIN`: [plain,http_simple,http_post,tls1.2_ticket_auth]",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("ss2ssr task: ", c.String("ss"), c.String("protocol"), c.String("obfs"))
				return nil
			},
		},
		{
			Name:     "list-ss",
			Category: "Converter",
			Usage:    "list shadowsocks from FROM[gui-config.json]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "gui-config.json",
					Usage: "specify `PATH` of gui-config.json",
				}},
			Action: func(c *cli.Context) error {
				fmt.Println("list-ss task: ", c.String("input"))
				return nil
			},
		},
		{
			Name:     "dnsmasq",
			Category: "Converter",
			Usage:    "list shadowsocks from FROM[gui-config.json]",
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
			Category: "Rank",
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
	app.Action = func(c *cli.Context) (err error) {
		SSRPATH = c.String("ssr")
		if _, err := os.Stat(SSRPATH); os.IsNotExist(err) {
			// path/to/whatever does not exist
			fmt.Println("SSR IsNotExist")
		}
		onlySSR = c.Bool("onlyssr")
		return
	}
	app.Run(os.Args)

	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	var (
		err error
	//	wg  sync.WaitGroup
	)

	if files := getFileList(getCurrDir(), isSSRList); files != nil {
		if guiCfgSSR, err = ParseConfig(files[0]); err != nil {
			panic(err)
		}
	} else {
		panic(errors.New(`NO ssr-list*.txt`))
	}

	if onlySSR {
		log.Println("Only Enable SSR", SSRPATH)
	} else {
		log.Println("Enable SS")
		guiCfgSS.Servers = make(SSRConfigSlice, len(guiCfgSSR.Servers))
		copy(guiCfgSS.Servers, guiCfgSSR.Servers)
	}

	go guiCfgSSR.TestServersBySSR()
	go guiCfgSS.TestServersBySS()
	go func() {
		if gfwList, err = ParseGFWList(guiCfgSS.GetServerArray()); err != nil {
			panic(err)
		}
	}()
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
