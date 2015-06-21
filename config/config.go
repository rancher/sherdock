package main

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)
type Config struct {
	pullIntervalMinutes int
	images []string
}

func configFileName(name string) string{
	if name == "" {
		name = "config.yaml"
	}
	return name
}

func defaultConfig() *Config{
	config := Config{}
	config.images = []string {}
	config.pullIntervalMinutes = 60
	return &config
}

func GetConfig(name string) (*Config, error) {
	a := Config{}
	b := true
	if name == "" {
		b = false
	}
	data, err := ioutil.ReadFile(configFileName(name))
	if err != nil && b {
		return  nil, err
	}
	if err != nil && !b {
		a := defaultConfig()
		SaveConfig(a, "")
		return a, nil
	}
	err = yaml.Unmarshal(data, &a)
	return &a , err
}

func SaveConfig(config *Config, name string) error {
	data, err := yaml.Marshal(&config)
	ioutil.WriteFile(configFileName(name), data, 0644)
	return err
}

func main() {
	config, err := GetConfig("")
	if  err == nil {
		fmt.Printf("%#v", config)
	} else  {
		fmt.Println(err)
	}
	config.images = append(config.images, "james")
	err = SaveConfig(config, "")
	if  err != nil {
		fmt.Println("You failed")
	}
}
