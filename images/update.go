package images

import (
	"github.com/samalba/dockerclient"
)

func PullImage(dockerClient *dockerclient.DockerClient) ([]*dockerclient.ImageInfo, error) {
	images, err := dockerClient.ListImages()
	if err != nil {
		return nil, err
	}
	var result = make([]*dockerclient.ImageInfo, len(images))
	for i, image := range images {
		tempResult, _ := InspectImage(dockerClient, image.Id)
		result[i] = tempResult
	}
	return result, nil
}
