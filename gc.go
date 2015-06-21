package main

import (
	"log"
	"github.com/samalba/dockerclient"
	"regexp"
)

func GC(dockerClient *dockerclient.DockerClient, regexFilter string) error {
	images, err := dockerClient.ListImages()
	if err != nil {
		return err
	}

	imagesToSave := make(map[string]bool)

	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			if match, _ := regexp.MatchString(regexFilter, repoTag); match {
				imagesToSave[image.Id] = true
			}
		}
	}

	for _, i := range images {
		if i.ParentId != "" {
			imagesToSave[i.ParentId] = true
		}
	}

	containers, err := dockerClient.ListContainers(true, false, "")
	if err != nil {
		return err
	}

	for _, c := range containers {
		info, _ := dockerClient.InspectContainer(c.Id)
		imagesToSave[info.Image] = true
	}


	for _, image := range images {
		if ! imagesToSave[image.Id] {
			log.Println("Deleting image with image id ", image.Id)
			dockerClient.RemoveImage(image.Id)
		}
	}
	return nil
}
