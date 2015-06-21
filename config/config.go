package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)
type Config struct {
	PullIntervalMinutes int
	Images []string
}

var conf = Config{}

func configFileName(name string) string{
	if name == "" {
		name = "config.yaml"
	}
	return name
}

func defaultConfig() *Config {
	config := Config{}
	config.Images = []string { "JAMESCONFIG" }
	config.PullIntervalMinutes = 60
	return &config
}

func GetConfig(name string) (*Config, error) {
	config := Config{}
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
	err = yaml.Unmarshal(data, &config)
	if  err == nil {
		return &config, nil
	} else  {
		return nil, err
	}
}

func SaveConfig(config *Config, name string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println("Failed to marshal the config.")
		return nil
	}
	fmt.Printf("%s", data)
	ioutil.WriteFile(configFileName(name), data, 0644)
	return err
}

func main1() {
	config, err := GetConfig("yaml")
	if err == nil {
		fmt.Printf( "Config %#v" , config)
	} else {
		fmt.Printf( "Err %#v", err)
	}
	conf = *config
	fmt.Printf("\n\n\n Config   %#v", conf)
}
