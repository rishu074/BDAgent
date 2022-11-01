package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type FtpStruct struct {
	Enabled bool   `yaml:"enabled"`
	FtpUrl  string `yaml:"uri"`
	User    string `yaml:"user"`
	Pass    string `yaml:"password"`
}

type ConfigStruct struct {
	Name          string    `yaml:"name"`
	Version       string    `yaml:"version"`
	Port          int       `yaml:"port"`
	Nodes         []string  `yaml:"nodes"`
	DataDirectory string    `yaml:"dataDirectory"`
	DataFileName  string    `yaml:"data_file"`
	Token         string    `yaml:"token"`
	BashFile      string    `yaml:"BashFile"`
	IpHeader      string    `yaml:"IP_HEADER"`
	Ftp           FtpStruct `yaml:"ftp"`
	ChunkSize     int       `yaml:"chunk_size"`
}

var data, _ = os.ReadFile("./config.yml")
var Conf = ConfigStruct{}
var _ = yaml.Unmarshal(data, &Conf)
