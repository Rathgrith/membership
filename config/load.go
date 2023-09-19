package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerList map[string]int `yaml:"server_list"`
	ServerPort int            `yaml:"server_port"`
	Introducer string         `yaml:"introducer"`
}

func LoadConfig() (Config, error) {
	yamlData, err := ioutil.ReadFile("config/host.yaml")
	if err != nil {
		return Config{}, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func GetHostID(hostname string) (int, error) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading:", err)
		return 0, err
	}
	id, exists := config.ServerList[hostname]
	if !exists {
		return 0, fmt.Errorf("hostname not found in mapping")
	}

	return id, nil
}

func GetIntroducer() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading:", err)
		return "", err
	}
	introducer := config.Introducer
	return introducer, nil
}
