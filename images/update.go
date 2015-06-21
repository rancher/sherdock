package images

import (
	"github.com/samalba/dockerclient"
	"fmt"
	"time"
	"github.com/rancherio/sherdock/config"
)

var url = "unix:///var/run/docker.sock"

func pullImages() error {
	dockerClient, err := dockerclient.NewDockerClient(url, nil)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}
	for {
		for _, element := range config.Conf.ImagesToPull {
			fmt.Printf("Pulling: %v\n", element)
			err := dockerClient.PullImage(element, nil)
			if err != nil {
				fmt.Printf("Error While Pulling: %#v\n", err)
			}
			fmt.Printf("Pulled: %v\n", element)
		}
		fmt.Println("All images pulled.\n")
		time.Sleep(time.Duration(config.Conf.GCIntervalMinutes) * time.Minute)
	}
	return  nil
}

func StartImageUpdate() {
	go pullImages()
}
