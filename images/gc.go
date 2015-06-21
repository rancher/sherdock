package images

import (
	"fmt"
	"regexp"
	"time"

	"github.com/rancherio/sherdock/config"
	"github.com/samalba/dockerclient"
)

func RunGC(docker *dockerclient.DockerClient, untagged bool, filters ...string) error {
	fmt.Println("Starting images GC")
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
				if del {
					fmt.Println("untagged ", image.Id)
				}
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
			fmt.Println("Deleteing image id ", image.Id, " name ", image.RepoTags)
			_, err := docker.RemoveImage(image.Id)
			if err != nil {
				fmt.Printf("Failed to delete %s: %v\n", image.Id, err)
			}
		}
	}

	fmt.Println("Done with images GC")

	return nil
}

func StartGC() error {
	for {
		client, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
		if err != nil {
			return err
		}

		RunGC(client, config.Conf.GCUntagged, config.Conf.ImagesToGC...)

		time.Sleep(time.Duration(config.Conf.GCIntervalMinutes) * time.Minute)
		config.LoadGlobalConfig()
	}
}
