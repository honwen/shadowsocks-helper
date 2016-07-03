package main

import (
	"fmt"
	"testing"
)

func TestConfigJson(t *testing.T) {
	config, err := ParseConfig("testdata/test-gui-config.json")
	if err != nil {
		t.Fatal("error parsing config.json:", err)
	}

	if len(config.Servers) != 3 {
		t.Error("wrong nums of Servers from config")
	}

	for i, server := range config.Servers {
		if server.Server != fmt.Sprintf("ss%02d.test.com", i+1) {
			t.Error("wrong Server of Server", i+1)
		}
		if server.Password != fmt.Sprintf("PPAASS%02d", i+1) {
			t.Error("wrong Password of Server", i+1)
		}
		if string(server.ServerPort) != fmt.Sprintf("%d", i+50001) {
			t.Error("wrong ServerPort of Server", i+1)
		}
		if server.Method != "rc4-md5" {
			t.Error("wrong Method of Server", i+1)
		}
	}
}
