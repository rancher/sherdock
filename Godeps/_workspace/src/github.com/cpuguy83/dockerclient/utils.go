package docker

import (
	"io"
	"strings"
)

func ParseURL(url string) (string, string) {
	arr := strings.Split(url, "://")

	if len(arr) == 1 {
		return "unix", arr[0]
	}

	proto := arr[0]
	if proto == "http" {
		proto = "tcp"
	}

	return proto, arr[1]
}

type readCloseWrapper struct {
	io.Reader
	closer func() error
}

func (r *readCloseWrapper) Close() error {
	return r.closer()
}

func newReadCloseWrapper(r io.Reader, closer func() error) io.ReadCloser {
	return &readCloseWrapper{
		Reader: r,
		closer: closer,
	}
}
