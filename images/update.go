package images

import (
	"github.com/samalba/dockerclient"
	"fmt"
)
var ToPull = []string { "tianon/true",  "rancher/server" }

var url = "unix:///var/run/docker.sock"

func PullImages() error {
	dockerClient, err := dockerclient.NewDockerClient(url, nil)
	if err != nil {
		return err
	}
	for _,element := range ToPull {
		err := dockerClient.PullImage(element, nil)
		if  err != nil {
			fmt.Printf("While Pulling: %#v", err)
		}
	}
	return  nil
}
