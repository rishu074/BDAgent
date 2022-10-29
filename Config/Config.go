package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Name          string   `yaml:"name"`
	Version       string   `yaml:"version"`
	Port          int      `yaml:"port"`
	Nodes         []string `yaml:"nodes"`
	DataDirectory string   `yaml:"dataDirectory"`
	DataFileName  string   `yaml:"data_file"`
	Token         string   `yaml:"token"`
	BashFile      string   `yaml:"BashFile"`
}

var data, _ = os.ReadFile("./config.yml")
var Conf = ConfigStruct{}
var _ = yaml.Unmarshal(data, &Conf)
