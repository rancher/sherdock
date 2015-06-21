package images

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"github.com/samalba/dockerclient"
)

var imageCache, _ = lru.New(1024)

func ListImagesDetailed(dockerClient *dockerclient.DockerClient) ([]*dockerclient.ImageInfo, error) {
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

func InspectImage(dockerClient *dockerclient.DockerClient, id string) (*dockerclient.ImageInfo, error) {
	cachedImageInfo, _ := imageCache.Get(id)
	if cachedImageInfo == nil {
		imageInfo, err := dockerClient.InspectImage(id)
		if err != nil {
			return nil, err
		}
		imageCache.Add(id, imageInfo)
		return imageInfo, nil
	} else {
		if cachedImageInfoCasted, ok := cachedImageInfo.(*dockerclient.ImageInfo); ok {
			return cachedImageInfoCasted, nil
		} else {
			return nil, fmt.Errorf("Cache casting error")
		}
	}
}
