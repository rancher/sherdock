package config

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	GCIntervalMinutes   int
	PullIntervalMinutes int
	ImagesToPull        []string
	ImagesToNotGC       []string
}

var Conf Config

func configFileName(name string) string {
	if name == "" {
		name = "config.yml"
	}
	return name
}

func defaultConfig() *Config {
	config := Config{
		GCIntervalMinutes:   5,
		PullIntervalMinutes: 60,
		ImagesToPull:        []string{"ubuntu:latest", "busybox:latest"},
		ImagesToNotGC:       []string{".*"},
	}

	return &config
}

func GetConfig(name string) (*Config, error) {
	config := Config{}
	data, err := ioutil.ReadFile(configFileName(name))
	if os.IsNotExist(err) {
		a := defaultConfig()
		SaveConfig(a, "")
		return a, nil
	}
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	return &config, err
}

func SaveConfig(config *Config, name string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println("Failed to marshal the config.")
		return err
	}
	ioutil.WriteFile(configFileName(name), data, 0644)
	return err
}

func LoadGlobalConfig() error {
	config, err := GetConfig(os.Getenv("SHERDOCK_CONFIG"))
	if err != nil {
		return err
	}
	Conf = *config

	return nil
}
