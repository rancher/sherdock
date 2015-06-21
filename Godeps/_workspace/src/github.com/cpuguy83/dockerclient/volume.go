package docker

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Volume struct {
	HostPath    string
	VolPath     string
	IsReadWrite bool
	IsBindMount bool
}

func (v *Volume) Id() string {
	return filepath.Base(v.HostPath)
}

func parseBindVolumeSpec(spec string) (Volume, error) {
	var (
		arr = strings.Split(spec, ":")
		vol Volume
	)

	switch len(arr) {
	case 1:
		vol.VolPath = spec
		vol.IsReadWrite = true
	case 2:
		vol.HostPath = arr[0]
		vol.VolPath = arr[1]
		vol.IsReadWrite = true
	case 3:
		vol.HostPath = arr[0]
		vol.VolPath = arr[1]
		vol.IsReadWrite = validVolumeMode(arr[2]) && arr[2] == "rw"
	default:
		return vol, fmt.Errorf("Invalid volume specification: %s", spec)
	}

	if !filepath.IsAbs(vol.HostPath) {
		return vol, fmt.Errorf("cannot bind mount volume: %s volume paths must be absolute.", vol.HostPath)
	}

	return vol, nil
}

func validVolumeMode(mode string) bool {
	validModes := map[string]bool{
		"rw": true,
		"ro": true,
	}

	return validModes[mode]
}
