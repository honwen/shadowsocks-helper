package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type SSConfig struct {
	Server     string `json:"server"`
	ServerPort string `json:"server_port"`
	// LocalPort  string `json:"local_port"`
	Password string `json:"password"`
	Method   string `json:"method"` // encryption method
	TestOK   bool
}

type Config struct {
	Servers []SSConfig `json:"configs"`
}

func (config *Config) GetServerArray() (Servers []string) {
	// ss://Method:Password@Server:ServerPort
	for _, ss := range config.Servers {
		Servers = append(Servers, fmt.Sprintf("TestOK: %t\t| ss://%s:%s@%s:%s",
			ss.TestOK, ss.Method, ss.Password, ss.Server, ss.ServerPort))
	}
	return
}

func (config *Config) TestServers() {
	var wg sync.WaitGroup
	for idx := range config.Servers {
		wg.Add(1)
		go func(idx int) {
			if _, err := wGetByShadowsocksProxy("http://google.com", config.Servers[idx].String()); err == nil {
				config.Servers[idx].TestOK = true
			}
			wg.Done()
		}(idx)
	}
	wg.Wait()
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

	config = &Config{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return
}

func (config Config) String() (str string) {
	for _, s := range config.GetServerArray() {
		str += s + "\n"
	}
	return
}

func (ss SSConfig) String() string {
	return fmt.Sprintf("ss://%s:%s@%s:%s", ss.Method, ss.Password, ss.Server, ss.ServerPort)
}
