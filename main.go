package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	guiCfg  *Config
	gfwList *GFWList
)

func main() {
	var (
		err error
		wg  sync.WaitGroup
	)
	if files := getFileList(getCurrDir(), isGuiConfig); files != nil {
		if guiCfg, err = ParseConfig(files[0]); err != nil {
			panic(err)
		}
	} else {
		panic(errors.New(`NO gui-config*.json`))
	}

	go guiCfg.TestServers()
	go func() {
		if gfwList, err = ParseGFWList(guiCfg.GetServerArray()); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		addr := ":1088"
		mux := http.NewServeMux()
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

		//===================================================================SS-RANK
		log.Println("Listen " + addr + "/ssrank")
		mux.HandleFunc("/ssrank", func(w http.ResponseWriter, r *http.Request) {
			if guiCfg.LastTestTime.Equal(NoLastTest) {
				io.WriteString(w, "// NOT READY")
				io.WriteString(w, "\n")
				log.Println("Listen " + addr + "/ssrank -- > NOT READY")
				return
			}
			if strings.HasPrefix(r.RequestURI, "/ssrank?") {
				switch r.RequestURI {
				case "/ssrank?renew": // renew crontrl, Duration > 1h
					if time.Now().Sub(guiCfg.LastTestTime) > 1*time.Hour {
						log.Println("Listen "+addr+"/ssrank", "accept renew")
						io.WriteString(w, fmt.Sprint("#", guiCfg.LastTestTime, ", ACCEPT", "\n"))
						go guiCfg.TestServers()
						guiCfg.LastTestTime = time.Now()
					} else {
						log.Println("Listen "+addr+"/ssrank", "reject renew")
						io.WriteString(w, fmt.Sprint("#", guiCfg.LastTestTime, ", REJECT", "\n"))
					}
				}
			} else {
				log.Println("Listen " + addr + "/ssrank")
				io.WriteString(w, fmt.Sprint("#", guiCfg.LastTestTime, "\n"))
				io.WriteString(w, guiCfg.GenTestedConfig().String())
			}
		})
		//===================================================================SS-RANK
		http.ListenAndServe(addr, mux)
	}()

	wg.Add(1)
	wg.Wait()
}
