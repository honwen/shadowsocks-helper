package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	TestCases = []Pair{
		Pair{"http://dl.google.com/pinyin/v2/GooglePinyinInstaller.exe", 60 * time.Second}, /*size: 15.7M, AtLeast 261KB/s*/
		Pair{"http://dl.google.com/dl/picasa/gpautobackup_setup.exe", 20 * time.Second},    /*size: 2.5M,  AtLeast 125KB/s*/
		Pair{"http://www.google.com/robots.txt", 3 * time.Second},                          /*size: 6.3K,  AtLeast 2.1KB/s*/
	}
	TestCaseIdx = 0
	NoLastTest  = time.Unix(0, 0)

	//isSSRForm = regexp.MustCompile(`ssr:\/\/([a-zA-Z0-9\.\-_]+:)(\d+:)([a-zA-Z0-9\.\-_]+:)+[a-zA-Z0-9]+`)
)

type Pair struct {
	a, b interface{}
}

type SSRConfig struct {
	Server     string      `json:"server"`
	ServerPort json.Number `json:"server_port"`
	Password   string      `json:"password"`
	Method     string      `json:"method"`
	Protocol   string      `json:"protocol"`
	Obfs       string      `json:"obfs"`
	Speed      float64     `json:"-"`
}

func (ssr SSRConfig) String() string {
	return fmt.Sprintf("ssr://%s:%s:%s:%s:%s:%s",
		ssr.Server, string(ssr.ServerPort), ssr.Protocol, ssr.Method, ssr.Obfs, ssr.Password)
}

func (ssr SSRConfig) ssString() string {
	return fmt.Sprintf("ss://%s:%s@%s:%s", ssr.Method, ssr.Password, ssr.Server, string(ssr.ServerPort))
}

type SSRConfigSlice []SSRConfig

func (s SSRConfigSlice) Len() int {
	return len(s)
}
func (s SSRConfigSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SSRConfigSlice) Less(i, j int) bool {
	return s[i].Speed < s[j].Speed
}

type Config struct {
	Servers      SSRConfigSlice
	LastTestTime time.Time
}

func (config *Config) GetServerArray() (Servers []string) {
	// ss://Method:Password@Server:ServerPort
	for _, ss := range config.Servers {
		Servers = append(Servers, fmt.Sprintf("ss://%s:%s@%s:%s",
			ss.Method, ss.Password, ss.Server, string(ss.ServerPort)))
	}
	return
}

func (config *Config) GenTestedConfig() (tested Config) {
	if config.LastTestTime.Equal(NoLastTest) {
		config.TestServersBySS()
	}

	tested.Servers = make(SSRConfigSlice, len(config.Servers))
	copy(tested.Servers, config.Servers)

	for idx, ss := range tested.Servers {
		if ss.Speed > 0 {
			continue
		}
		tested.Servers = tested.Servers[:idx] // Cut 1st Bad
		return
	}
	return
}

func (config *Config) GetBestServerArray() (Servers []string) {
	// ss://Method:Password@Server:ServerPort
	tested := config.GenTestedConfig()
	return tested.GetServerArray()
}

func (config *Config) TestServersBySS() {
	var wg sync.WaitGroup
	for idx := range config.Servers {
		wg.Add(1)
		go func(idx int) {
			log.Println(fmt.Sprintf("%02d", idx), "Test begin:", config.Servers[idx].ssString())
			tsBegin := float64(time.Now().UnixNano()) / 1000 // timestaps of Microsecond
			if bytes, err := wGetRawFastByShadowsocksProxy(
				TestCases[TestCaseIdx].a.(string),
				config.Servers[idx].ssString(),
				TestCases[TestCaseIdx].b.(time.Duration),
			); err == nil {
				tsDone := float64(time.Now().UnixNano()) / 1000                      // timestaps of Microsecond
				config.Servers[idx].Speed = float64(len(bytes)) / (tsDone - tsBegin) // Bytes / Microsecond => MB/s
				log.Println(fmt.Sprintf("%02d", idx), "Test done:", float64(len(bytes))/1000, "Kbytes in", (tsDone-tsBegin)/1000000, "seconds")
			} else {
				config.Servers[idx].Speed = -1
				log.Println(fmt.Sprintf("%02d", idx), "Test done:", "Fail of", err)
			}
			wg.Done()
		}(idx)
		if idx > 0 && idx%3 == 0 {
			time.Sleep(10 * time.Second)
			// Avoid too much connection in the same time
		} else {
			time.Sleep(888 * time.Millisecond)
		}
	}
	wg.Wait()
	if !sort.IsSorted(sort.Reverse(config.Servers)) {
		sort.Sort(sort.Reverse(config.Servers))
	}
	config.LastTestTime = time.Now()
}

func (config *Config) TestServersBySSR() {
	var wg sync.WaitGroup
	for idx := range config.Servers {
		wg.Add(1)
		go func(idx int) {
			log.Println(fmt.Sprintf("%02d", idx), "Test begin:", config.Servers[idx].String())
			tsBegin := float64(time.Now().UnixNano()) / 1000 // timestaps of Microsecond
			if bytes, err := wGetRawFastByShadowsocksRProxy(
				"",
				TestCases[TestCaseIdx].a.(string),
				config.Servers[idx].String(),
				TestCases[TestCaseIdx].b.(time.Duration),
			); err == nil {
				tsDone := float64(time.Now().UnixNano()) / 1000                                    // timestaps of Microsecond
				config.Servers[idx].Speed = float64(len(bytes)) / (tsDone - tsBegin - 3*1000*1000) // Bytes / Microsecond => MB/s
				log.Println(fmt.Sprintf("%02d", idx), "Test done:", float64(len(bytes))/1000, "Kbytes in", (tsDone-tsBegin)/1000000, "seconds")
			} else {
				config.Servers[idx].Speed = -1
				log.Println(fmt.Sprintf("%02d", idx), "Test done:", "Fail of", err)
			}
			wg.Done()
		}(idx)
		if idx > 0 && idx%3 == 0 {
			time.Sleep(10 * time.Second)
			// Avoid too much connection in the same time
		} else {
			time.Sleep(888 * time.Millisecond)
		}
	}
	wg.Wait()
	if !sort.IsSorted(sort.Reverse(config.Servers)) {
		sort.Sort(sort.Reverse(config.Servers))
	}
	config.LastTestTime = time.Now()
}

func ParseConfig(path string) (config *Config, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return
	}
	defer file.Close()

	config = &Config{
		Servers:      make(SSRConfigSlice, 0),
		LastTestTime: NoLastTest,
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	ssr := isSSRForm.FindAllString(string(data), -1)

	for _, s := range ssr {
		// fmt.Println(s)
		if strings.Count(s, `:`) >= 5 {
			s = strings.TrimPrefix(s, `ssr://`)
			s = strings.TrimSuffix(s, "\n")
			s = strings.TrimSuffix(s, "\r")
			sub := strings.Split(s, ":")
			config.Servers = append(config.Servers, SSRConfig{
				Server:     sub[0],
				ServerPort: json.Number(sub[1]),
				Protocol:   sub[2],
				Method:     sub[3],
				Obfs:       sub[4],
				Password:   sub[5],
			})
		} else {
			// TODO add full base64 SSR-QRcode-Scheme Support
		}

	}

	// fmt.Println(config.Servers)
	return
}

func (config Config) String() (str string) {
	for _, ssr := range config.Servers {
		str += fmt.Sprintf("Speed:%6.03fMB/s, ssr://%s:%s:%s:%s:%s:%s\n",
			ssr.Speed, ssr.Server, string(ssr.ServerPort), ssr.Protocol, ssr.Method, ssr.Obfs, ssr.Password)
	}
	return
}

func (config Config) ssString() (str string) {
	for _, ss := range config.Servers {
		str += fmt.Sprintf("Speed:%6.03fMB/s, ss://%s:%s@%s:%s\n",
			ss.Speed, ss.Method, ss.Password, ss.Server, string(ss.ServerPort))
	}
	return
}
