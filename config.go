package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	TestCases = []Pair{
		Pair{"http://dl.google.com/pinyin/v2/GooglePinyinInstaller.exe", 60 * time.Second}, /*size: 15.7M, AtLeast 261KB/s*/
		Pair{"http://dl.google.com/dl/picasa/gpautobackup_setup.exe", 20 * time.Second},    /*size: 2.5M,  AtLeast 125KB/s*/
		Pair{"http://www.google.com/robots.txt", 3 * time.Second},                          /*size: 6.3K,  AtLeast 2.1KB/s*/
	}
	TestCaseIdx = 1
	NoLastTest  = time.Unix(0, 0)
)

type Pair struct {
	a, b interface{}
}

type SSConfig struct {
	Server     string      `json:"server"`
	ServerPort json.Number `json:"server_port"`
	// LocalPort  json.Number `json:"local_port"`
	Password string `json:"password"`
	Method   string `json:"method"` // encryption method
	Speed    float64
}

func (ss SSConfig) String() string {
	return fmt.Sprintf("ss://%s:%s@%s:%s", ss.Method, ss.Password, ss.Server, string(ss.ServerPort))
}

type SSConfigSlice []SSConfig

func (s SSConfigSlice) Len() int {
	return len(s)
}
func (s SSConfigSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SSConfigSlice) Less(i, j int) bool {
	return s[i].Speed < s[j].Speed
}

type Config struct {
	Servers      SSConfigSlice `json:"configs"`
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
		config.TestServers()
	}
	tested = *config
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

func (config *Config) TestServers() {
	var wg sync.WaitGroup
	for idx := range config.Servers {
		wg.Add(1)
		go func(idx int) {
			log.Println(fmt.Sprintf("%02d", idx), "Test begin:", config.Servers[idx].String())
			tsBegin := float64(time.Now().UnixNano()) / 1000 // timestaps of Microsecond
			if bytes, err := wGetRawFastByShadowsocksProxy(
				TestCases[TestCaseIdx].a.(string),
				config.Servers[idx].String(),
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
		if idx > 0 && idx%4 == 0 {
			time.Sleep(4 * time.Second)
			// Avoid to much connection in the same time
		} else {
			time.Sleep(444 * time.Millisecond)
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

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	config = &Config{
		LastTestTime: NoLastTest,
	}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return
}

func (config Config) String() (str string) {
	for _, ss := range config.Servers {
		str += fmt.Sprintf("Speed:%6.03fMB/s, ss://%s:%s@%s:%s\n",
			ss.Speed, ss.Method, ss.Password, ss.Server, string(ss.ServerPort))
	}
	return
}

func ss2json(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil || u.User == nil {
		return ``, err
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return ``, err
	}
	ssJSON := struct {
		Server     string      `json:"server"`
		ServerPort json.Number `json:"server_port"`
		Password   string      `json:"password"`
		Method     string      `json:"method"`
	}{
		Server:     host,
		ServerPort: json.Number(port),
		Method:     u.User.Username(),
	}
	ssJSON.Password, _ = u.User.Password()
	b, err := json.MarshalIndent(ssJSON, "", "  ")
	if err != nil {
		return ``, err
	}
	return string(b), nil
}
