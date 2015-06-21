package images

import (
	"github.com/samalba/dockerclient"
	"fmt"
	"time"
	"github.com/rancherio/sherdock/config"
)
var ToPull = []string { "tianon/true",  "rancher/server" }

var url = "unix:///var/run/docker.sock"
var PullDelay = 1

func pullImages() error {
	fmt.Println("Making clinet")
	dockerClient, err := dockerclient.NewDockerClient(url, nil)
	if err != nil {
		fmt.Printf("%#v", err)
		return err
	}
	for {
		fmt.Println("%#v", config.Conf)
		for _, element := range config.Conf.ImagesToPull {
			fmt.Printf("Pulling: %v", element)
			err := dockerClient.PullImage(element, nil)
			if err != nil {
				fmt.Printf("Error While Pulling: %#v", err)
			}
			fmt.Printf("Pulled: %v", element)
		}
		fmt.Println("All images pulled.")
		time.Sleep(time.Duration(config.Conf.GCIntervalMinutes) * time.Minute)
	}
	return  nil
}

func StartImageUpdate() {
	go pullImages()
}
