package main

import (
	"github.com/samalba/dockerclient"
)

func ListContainersDetailed(dockerClient *dockerclient.DockerClient) ([]*dockerclient.ContainerInfo, error) {
	containers, err := dockerClient.ListContainers(true, false, "")
	if err != nil {
		return nil, err
	}
	var result = make([]*dockerclient.ContainerInfo,len(containers))
	for i, container := range containers {
		containerInfo, err := dockerClient.InspectContainer(container.Id)
		if err != nil {
			return nil, err
		}
		result[i] = containerInfo
	}
	return result, nil
}