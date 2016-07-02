package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type SSConfig struct {
	Server     string `json:"server"`
	ServerPort string `json:"server_port"`
	// LocalPort  string `json:"local_port"`
	Password string `json:"password"`
	Method   string `json:"method"` // encryption method
}

type Config struct {
	Servers []SSConfig `json:"configs"`
}

func (config *Config) GetServerArray() (Servers []string) {
	// ss://Method:Password@Server:ServerPort
	for _, ss := range config.Servers {
		Servers = append(Servers, fmt.Sprintf("ss://%s:%s@%s:%s",
			ss.Method, ss.Password, ss.Server, ss.ServerPort))
	}
	return
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
