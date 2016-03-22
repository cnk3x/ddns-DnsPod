package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	LogConfig LogConfig      `json:"log"`
	Records   []RecordConfig `json:"records"`
}

type LogConfig struct {
	Debug bool `json:debug`
	Info  bool `json:info`
	Error bool `json:error`
}

type RecordConfig struct {
	Token  string `json:"token"`
	Host   string `json:"host"`
	Domain string `json:"domain"`
}

var config Config

func init() {
	loadConfig()
}

func loadConfig() error {
	body, err := ioutil.ReadFile(getConfigFile())
	err = json.Unmarshal(body, &config)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFile() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)

	paths := strings.Split(path, p)
	fileName := p + paths[len(paths)-1]
	logFile := p + "ddns.config"
	path = strings.Replace(path, fileName, logFile, 1)
	return path
}
