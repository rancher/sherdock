package main

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)
type Config struct {
	PullIntervalMinutes int
	Images []string
}

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
		fmt.Println( "Failing" )
		return  nil, err
	}
	if err != nil && !b {
		a := defaultConfig()
		SaveConfig(a, "")
		return a, nil
	}
	fmt.Printf( "%s", data )
	err = yaml.Unmarshal(data, &config)
	fmt.Printf( "%#v After unmarshal \n", config)
	if  err == nil {
		fmt.Printf( "%#v No Err\n", config)
		return &config, nil
	} else  {
		fmt.Printf( "%#v Err\n", config)
		fmt.Printf( "%v Err\n", err)
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

func main() {
	config, err := GetConfig("")
	if err == nil {
		fmt.Printf( "Config %#v" , config)
	} else {
		fmt.Printf( "Err %#v", err)
	}
}
