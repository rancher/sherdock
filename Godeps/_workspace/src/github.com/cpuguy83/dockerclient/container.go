package docker

import (
	"fmt"
	"path/filepath"
)

type Container struct {
	Id              string
	Name            string
	NetworkSettings *NetworkSettings
	State           struct {
		Running bool
	}
	Config struct {
		Image        string
		AttachStderr bool
		AttachStdin  bool
		AttachStdout bool
	}
	HostConfig struct {
		PortBindings map[string][]Binding
		Binds        []string
	}
	Volumes   map[string]string
	VolumesRW map[string]bool
}

func (container *Container) GetVolumes() (map[string]*Volume, error) {
	// Get all the bind-mounts
	volumes, err := container.getBindMap()
	if err != nil {
		return nil, err
	}

	// Get all the normal volumes
	for volPath, hostPath := range container.Volumes {
		if _, exists := volumes[volPath]; exists {
			continue
		}
		volumes[volPath] = &Volume{VolPath: volPath, HostPath: hostPath, IsReadWrite: container.VolumesRW[volPath]}
	}

	return volumes, nil
}

func (container *Container) getBindMap() (map[string]*Volume, error) {
	var (
		// Create the requested bind mounts
		volumes = map[string]*Volume{}
		// Define illegal container destinations
		illegalDsts = []string{"/", "."}
	)

	for _, bind := range container.HostConfig.Binds {
		vol, err := parseBindVolumeSpec(bind)
		if err != nil {
			return nil, err
		}
		vol.IsBindMount = true
		// Bail if trying to mount to an illegal destination
		for _, illegal := range illegalDsts {
			if vol.VolPath == illegal {
				return nil, fmt.Errorf("Illegal bind destination: %s", vol.VolPath)
			}
		}

		volumes[filepath.Clean(vol.VolPath)] = &vol
	}
	return volumes, nil
}
