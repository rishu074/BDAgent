package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Name    string
	Version string
	Port    int
}

var data, _ = os.ReadFile("./config.yml")
var Conf = ConfigStruct{}
var _ = yaml.Unmarshal(data, &Conf)
