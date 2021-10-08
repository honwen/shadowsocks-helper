package main

import (
	"bytes"
	"fmt"
	"io"
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
	"github.com/Workiva/go-datastructures/set"
	"github.com/chenhw2/go-ps" /*ps*/
	"github.com/chenhw2/shadowsocks-helper/cidr"
	"github.com/chenhw2/shadowsocks-helper/ssStruct"
	"github.com/chenhw2/shadowsocks-helper/subscribe"
	"github.com/honwen/golibs/cip"
	"github.com/honwen/golibs/domain"
	"gopkg.in/urfave/cli.v1"
)

const ipipCHN = "https://raw.githubusercontent.com/17mon/china_ip_list/master/china_ip_list.txt"

var (
	version        = "MISSING build version [git hash]"
	defaultSSRPath = "ssr-local"
	defaultASNs    = cli.Int64Slice{cip.ASN_TENCENT_CN, cip.ASN_ALIBABA_CN, cip.ASN_HINET, cip.ASN_HKBN, cip.ASN_NEWTT, cip.ASN_UHGL, cip.ASN_MICROSOFT, cip.ASN_BCPL_SG}
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
				}},
			Action: func(c *cli.Context) error {
				// fmt.Println("dnsmasq task: ", c.String("server"), c.String("ipset"), c.String("proxy"), c.String("path"))
				domains, err := gfwlist()

				if err != nil {
					log.Printf("%+v", err)
				} else {
					var (
						str bytes.Buffer
						w   io.Writer = &str
					)
					for i := range domains {
						io.WriteString(w, fmt.Sprintf("server=/%s/%s\n", domains[i], c.String("server")))
					}
					if !strings.EqualFold(c.String("ipset"), "NULL") {
						for i := range domains {
							io.WriteString(w, fmt.Sprintf("ipset=/%s/%s\n", domains[i], c.String("ipset")))
						}
					}
					fmt.Println(str.String())
				}
				return nil
			},
		},
		{
			Name:     "gfwlist",
			Category: "HELPER",
			Usage:    "generate Domain-List from GFWLIST(online) with SSR-Proxy(optional)",
			Action: func(c *cli.Context) error {
				// fmt.Println("gfwlist task: ", c.String("proxy"), c.String("path"))
				domains, err := gfwlist()
				if err != nil {
					log.Printf("%+v", err)
				} else {
					gfwlist := strings.Join(domains, "\n")
					fmt.Println(gfwlist)
				}
				return nil
			},
		},
		{
			Name:     "tide",
			Category: "HELPER",
			Usage:    "tide Domain-List",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "path of domain list",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "path of domain list",
				}},
			Action: func(c *cli.Context) error {
				// fmt.Println("tide task: ", c.String("input"), c.String("output"))

				if dat, err := ioutil.ReadFile(c.String("input")); err != nil {
					log.Printf("%+v", err)
					return nil
				} else {
					tidelist := strings.Join(domain.Sort(domain.ExtractFromBytes(dat)), "\n")
					if err != nil {
						log.Printf("%+v", err)
					} else {
						if len(c.String("output")) > 0 {
							err := ioutil.WriteFile(c.String("output"), []byte(tidelist), 0644)
							if err != nil {
								panic(err)
							}
						} else {
							fmt.Println(tidelist)
						}
					}
				}
				return nil
			},
		},
		{
			Name:     "asn",
			Category: "HELPER",
			Usage:    "get IPs of ASN",
			Flags: []cli.Flag{
				cli.Int64SliceFlag{
					Name:  "asn, a",
					Usage: "Number of ASN",
					Value: &defaultASNs,
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "path of domain list",
				}},
			Action: func(c *cli.Context) error {
				// fmt.Println("asn task: ", c.Int64Slice("asn"), len(c.Int64Slice("asn")), c.String("output"))
				asn := c.Int64Slice("asn")
				ips := set.New()
				content := ""

				rlen := len(asn)
				rchan := make(chan []string, rlen)
				for _, it := range asn {
					go func(it int64) {
						rchan <- cip.IPsOfASN(int(it))
					}(it)
				}

				for i := 0; i < rlen; i++ {
					for _, v := range <-rchan {
						ips.Add(v)
					}
				}

				for _, it := range ips.Flatten() {
					content += fmt.Sprintln(it)
				}
				if len(c.String("output")) > 0 {
					err := ioutil.WriteFile(c.String("output"), []byte(content), 0644)
					if err != nil {
						panic(err)
					}
				} else {
					fmt.Println(content)
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

func gfwlist() (domains []string, err error) {
	domains = append(domains, domain.ExtractFromBytes([]byte(initList))...)

	var (
		topDomains = []string{}
		b64len     = len(officalGFWListURLs)
		comlen     = len(communityDomainLists)
		b64chan    = make(chan string, b64len)
		comchan    = make(chan string, comlen)

		toplen  = len(topDomainLists)
		topchan = make(chan string, toplen)
	)

	for i := 0; i < b64len; i++ {
		go func(it int) {
			if base64Str, err := wGet(officalGFWListURLs[it]); err == nil {
				b64chan <- base64Str
			} else {
				log.Println(err)
				b64chan <- ""
			}
		}(i)
	}

	for i := 0; i < comlen; i++ {
		go func(it int) {
			if extraList, err := wGet(communityDomainLists[it]); err == nil {
				comchan <- extraList
			} else {
				log.Println(err)
				comchan <- ""
			}
		}(i)
	}

	for i := 0; i < toplen; i++ {
		go func(it int) {
			if extraList, err := wGet(topDomainLists[it]); err == nil {
				topchan <- extraList
			} else {
				log.Println(err)
				topchan <- ""
			}
		}(i)
	}

	for i := 0; i < b64len; i++ {
		base64Str := <-b64chan
		domains = append(domains, domain.ExtractFromB64(base64Str)...)
	}

	for i := 0; i < comlen; i++ {
		extraList := <-comchan
		domains = append(domains, domain.ExtractFromBytes([]byte(extraList))...)
	}

	for i := 0; i < toplen; i++ {
		extraList := <-topchan
		topDomains = append(topDomains, domain.ExtractFromBytes([]byte(extraList))...)
	}

	_domains := []string{}
	for _, v := range domains {
		_domains = append(_domains, strings.TrimPrefix(v, "www."))
	}

	domains = customsSort(_domains, topDomains)
	return
}

// Get return text of url
func wGet(urlAddr string) (string, error) {
	request, err := http.NewRequest("GET", urlAddr, nil)
	if err != nil {
		return ``, err
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
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
