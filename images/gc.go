package images

import (
	"fmt"
	"regexp"

	"github.com/samalba/dockerclient"
)

func RunGC(docker *dockerclient.DockerClient, untagged bool, filters ...string) error {
	// list all the containers
	containers, err := docker.ListContainers(true, false, "")
	if err != nil {
		return err
	}

	hasChild := make(map[string]bool)
	inUse := make(map[string]bool)

	images, err := docker.ListImages()
	if err != nil {
		return err
	}

	for _, i := range images {
		hasChild[i.ParentId] = true
	}

	for _, c := range containers {
		info, _ := docker.InspectContainer(c.Id)
		inUse[info.Image] = true
	}

	for _, image := range images {
		if hasChild[image.Id] || inUse[image.Id] {
			continue
		}

		del := false
	outer:
		for _, tag := range image.RepoTags {
			if tag == "<none>:<none>" {
				del = untagged
				break
			}

			for _, filter := range filters {
				matched, err := regexp.Match(filter, []byte(tag))
				if err != nil {
					return err
				}
				if matched {
					fmt.Println("matched ", tag, " to ", filter)
					del = true
					break outer
				}
			}

		}

		if del {
			fmt.Println("Found image id ", image.Id, " name ", image.RepoTags)
			//_, err := docker.RemoveImage(image)
			//if err != nil {
			//log.Error("Failed to delete %s: %v", image.Id, err)
			//}
		}
	}

	return nil
}
