package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
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
	myApp := cli.NewApp()
	myApp.Name = "ss-helper"
	myApp.Usage = "shadowsocks(R)-helper"
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ssr",
			Value: "ssr-local",
			Usage: "Set SSR path",
		},
		cli.BoolFlag{
			Name:  "onlyssr",
			Usage: "Enable SS",
		},
	}
	myApp.Action = func(c *cli.Context) (err error) {
		SSRPATH = c.String("ssr")
		if _, err := os.Stat(SSRPATH); os.IsNotExist(err) {
			// path/to/whatever does not exist
			fmt.Println("SSR IsNotExist")
		}
		onlySSR = c.Bool("onlyssr")
		return
	}
	myApp.Run(os.Args)
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	var (
		err error
		wg  sync.WaitGroup
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
}
