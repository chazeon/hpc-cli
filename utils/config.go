package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Machine struct {
	Name string `json:"name"`
	Host string `json:"host"`
	User string `json:"user"`
}

type Config struct {
	Machines []Machine
	Commands map[string]string
	AuthKey  string `yaml:"auth_key"`
}

func LoadConfig(fname string) (config Config, err error) {

	bytes, err := os.ReadFile(fname)

	if err != nil {
		return
	}

	err = yaml.Unmarshal(bytes, &config)

	return
}
