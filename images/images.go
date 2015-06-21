package images

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"github.com/samalba/dockerclient"
)

type DetailedImageInfo struct {
	Created     int64
	Id          string
	ParentId    string
	RepoTags    []string
	Size        int64
	VirtualSize int64

	Architecture    string
	Author          string
	Comment         string
	//Config          *ContainerConfig
	Container       string
	//ContainerConfig *ContainerConfig
	DockerVersion   string
	Os              string
}

var imageCache, _ = lru.New(1024)

func ListImagesDetailed(dockerClient *dockerclient.DockerClient) ([]*DetailedImageInfo, error) {
	images, err := dockerClient.ListImages()
	if err != nil {
		return nil, err
	}
	var result = make([]*DetailedImageInfo, len(images))
	for i, image := range images {
		imagesDetails, _ := InspectImage(dockerClient, image.Id)
		detailedImageInfo := DetailedImageInfo{
			Created: image.Created,
			Id: image.Id,
			ParentId: image.ParentId,
			RepoTags: image.RepoTags,
			Size: image.Size,
			VirtualSize: image.VirtualSize,
			Architecture: imagesDetails.Architecture,
			Author: imagesDetails.Author,
			Comment: imagesDetails.Comment,
			Container: imagesDetails.Container,
			DockerVersion: imagesDetails.Container,
			Os: imagesDetails.Os }
		result[i] = &detailedImageInfo
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
