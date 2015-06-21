package docker

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type (
	Docker interface {
		FetchAllContainers(all bool) ([]*Container, error)
		FetchContainer(name string) (*Container, error)
		GetEvents() chan *Event
		Info() (*DaemonInfo, error)
		PullImage(name string) error
		CreateContainer(container map[string]interface{}) (string, error)
		StartContainer(string, interface{}) error
		RunContainer(map[string]interface{}) (string, error)
		RemoveContainer(name string, force, volumes bool) error
		ContainerLogs(id string, follow, stdout, stderr, timestamps bool, tail int) (io.ReadCloser, error)
		ContainerPause(id string) error
		ContainerUnpause(id string) error
		Copy(id string, file string) (io.ReadCloser, error)
		Build(ctx io.Reader, tag string, nocache bool, forcerm bool) (io.ReadCloser, error)
		DecodeStream(stream io.Reader) []string
		RemoveImage(name string, force bool, noprune bool) (io.ReadCloser, error)
		ContainerWait(name string) error
		SetTlsConfig(config *tls.Config)
		Version() (*DaemonVersion, error)
		ContainerStats(name string) (io.ReadCloser, error)
		//Attach(name string, logs, stream, stdin, stdout, stderr bool) (io.Reader, io.Writer, error)
	}

	Event struct {
		ContainerId string `json:"id"`
		Status      string `json:"status"`
	}

	Binding struct {
		HostIp   string
		HostPort string
	}

	NetworkSettings struct {
		IpAddress string
		Ports     map[string][]Binding
	}

	dockerClient struct {
		path      string
		tlsConfig *tls.Config
	}

	DaemonInfo struct {
		Containers         int
		Images             int
		Driver             string
		DriverStatus       [][]string
		ExecutionDriver    string
		KernelVersion      string
		NCPU               int
		MemTotal           int64
		Name               string
		ID                 string
		Debug              int
		NFd                int
		NGoroutines        int
		NEventsListener    int
		InitPath           string
		InitSha1           string
		IndexServerAddress string
		MemoryLimit        int
		SwapLimit          int
		IPv4Forwarding     int
		Labels             []string
		DockerRootDir      string
		OperatingSystem    string
	}

	DaemonVersion struct {
		ApiVersion    string
		Arch          string
		GitCommit     string
		GoVersion     string
		KernelVersion string
		Os            string
		Version       string
	}
)

func (d *DaemonInfo) RootPath() string {
	for _, i := range d.DriverStatus {
		if i[0] == "Root Dir" {
			return i[1]
		}
	}
	return ""
}

func NewClient(path string) (Docker, error) {
	return &dockerClient{path: path}, nil
}

func (d *dockerClient) SetTlsConfig(config *tls.Config) {
	d.tlsConfig = config
}

func (d *dockerClient) newConn() (*httputil.ClientConn, error) {
	var (
		conn net.Conn
		err  error
	)
	proto, path := ParseURL(d.path)
	if d.tlsConfig == nil {
		conn, err = net.Dial(proto, path)
	} else {
		conn, err = tls.Dial(proto, path, d.tlsConfig)
	}

	if err != nil {
		return nil, err
	}
	return httputil.NewClientConn(conn, nil), nil
}

func (docker *dockerClient) PullImage(name string) error {
	var (
		method = "POST"
		uri    = fmt.Sprintf("/images/create?fromImage=%s", name)
	)

	respBody, err := docker.newRequest(method, uri, nil)
	if err != nil {
		return nil
	}

	respBody.Close()

	return nil
}

func (docker *dockerClient) RemoveContainer(name string, force, volumes bool) error {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("/containers/%s?force=%s&volumes=%s", name, strconv.FormatBool(force), strconv.FormatBool(volumes))
	)

	respBody, err := docker.newRequest(method, uri, nil)
	if err != nil {
		return err
	}
	respBody.Close()

	return nil
}

func (docker *dockerClient) CreateContainer(container map[string]interface{}) (string, error) {
	var (
		method = "POST"
		name   string
		uri    = "/containers/create"
	)

	if _, exists := container["Name"]; exists {
		uri = fmt.Sprintf("%s?name=%v", uri, name)
	}

	delete(container, "Name")
	respBody, err := docker.newRequest(method, uri, container)
	if err != nil {
		// Try to see if we just need to download the image
		if fmt.Sprintf("%v", err) == "invalid HTTP request 404 404 Not Found" {
			if err := docker.PullImage(fmt.Sprintf("%s", container["Image"])); err != nil {
				return "", err
			}
			respBody, err = docker.newRequest(method, uri, container)
		}
		if err != nil {
			return "", err
		}
	}
	defer respBody.Close()

	type createResp struct {
		Id string
	}
	var respData createResp
	err = json.NewDecoder(respBody).Decode(&respData)
	if err != nil {
		return name, err
	}
	name = respData.Id

	return name, nil
}

func (docker *dockerClient) StartContainer(name string, hostConfig interface{}) error {
	var (
		method = "POST"
		uri    = fmt.Sprintf("/containers/%s/start", name)
	)

	respBody, err := docker.newRequest(method, uri, hostConfig)
	if err != nil {
		return err
	}
	defer respBody.Close()

	return nil
}

func (docker *dockerClient) RunContainer(config map[string]interface{}) (string, error) {

	name, err := docker.CreateContainer(config)
	if err != nil {
		return "", err
	}

	return name, docker.StartContainer(name, config["HostConfig"])
}

func (docker *dockerClient) FetchContainer(name string) (*Container, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("/containers/%s/json", name)
	)

	respBody, err := docker.newRequest(method, uri, nil)

	if err != nil {
		return nil, err
	}
	defer respBody.Close()
	var container *Container
	err = json.NewDecoder(respBody).Decode(&container)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func (docker *dockerClient) FetchAllContainers(all bool) ([]*Container, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("/containers/json?all=%v", all)
	)

	respBody, err := docker.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}
	defer respBody.Close()

	var containers []*Container
	if err = json.NewDecoder(respBody).Decode(&containers); err != nil {
		return nil, err
	}
	return containers, nil
}

func (docker *dockerClient) newRequest(method, uri string, body interface{}) (io.ReadCloser, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	c, err := docker.newConn()
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if !docker.isOkStatus(resp.StatusCode) {
		return nil, fmt.Errorf("invalid HTTP request %d %s", resp.StatusCode, resp.Status)
	}

	r := newReadCloseWrapper(resp.Body, func() error {
		resp.Body.Close()
		return c.Close()
	})

	return r, nil
}

func (d *dockerClient) isOkStatus(code int) bool {
	codes := map[int]bool{
		200: true,
		201: true,
		204: true,
		400: false,
		404: false,
		500: false,
		409: false,
		406: false,
	}

	return codes[code]
}

func (docker *dockerClient) Info() (*DaemonInfo, error) {
	var (
		method = "GET"
		uri    = "/info"
	)

	respBody, err := docker.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}
	defer respBody.Close()

	var info *DaemonInfo
	if err = json.NewDecoder(respBody).Decode(&info); err != nil {
		return nil, err
	}
	return info, nil
}

func (docker *dockerClient) Version() (*DaemonVersion, error) {
	var (
		method = "GET"
		uri    = "/version"
	)

	respBody, err := docker.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}
	defer respBody.Close()

	var version *DaemonVersion
	if err = json.NewDecoder(respBody).Decode(&version); err != nil {
		return nil, err
	}
	return version, nil
}

func (d *dockerClient) GetEvents() chan *Event {
	eventChan := make(chan *Event, 100) // 100 event buffer
	go func() {
		defer close(eventChan)

		respBody, err := d.newRequest("GET", "/events", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer respBody.Close()

		// handle signals to stop the socket
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			for sig := range sigChan {
				log.Printf("received signal '%v', exiting", sig)

				respBody.Close()
				close(eventChan)
				os.Exit(0)
			}
		}()

		dec := json.NewDecoder(respBody)
		for {
			var event *Event
			if err := dec.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("cannot decode json: %s", err)
				continue
			}
			eventChan <- event
		}
	}()
	return eventChan
}

func (d *dockerClient) ContainerLogs(id string, follow, stdout, stderr, timestamps bool, tail int) (io.ReadCloser, error) {
	tailStr := strconv.Itoa(tail)
	if tail == -1 {
		tailStr = "all"
	}
	uri := fmt.Sprintf("/containers/%s/logs?follow=%v&stdout=%v&stderr=%v&timestamps=%v&tail=%v", id, follow, stdout, stderr, timestamps, tailStr)

	respBody, err := d.newRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (d *dockerClient) Copy(id string, file string) (io.ReadCloser, error) {
	var (
		method = "POST"
		uri    = fmt.Sprintf("/containers/%s/copy", id)
		body   = map[string]string{"Resource": file}
	)

	respBody, err := d.newRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (d *dockerClient) ContainerPause(id string) error {
	var (
		method = "POST"
		uri    = fmt.Sprintf("/containers/%s/pause", id)
	)
	respBody, err := d.newRequest(method, uri, nil)
	if err != nil {
		return err
	}
	respBody.Close()
	return nil
}

func (d *dockerClient) ContainerUnpause(id string) error {
	var (
		method = "POST"
		uri    = fmt.Sprintf("/containers/%s/unpause", id)
	)
	respBody, err := d.newRequest(method, uri, nil)
	if err != nil {
		return err
	}
	respBody.Close()

	return nil
}

func (d *dockerClient) Build(ctx io.Reader, tag string, nocache bool, forcerm bool) (io.ReadCloser, error) {
	var (
		method = "POST"
		uri    = "/build"
	)

	v := &url.Values{}
	if tag != "" {
		v.Set("t", tag)
	}
	if nocache {
		v.Set("nocache", "1")
	}
	v.Set("rm", "1")
	if forcerm {
		v.Set("forcerm", "1")
	}

	uri = fmt.Sprintf("%s?%s", uri, v.Encode())
	req, err := http.NewRequest(method, uri, ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/tar")

	c, err := d.newConn()
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if !d.isOkStatus(resp.StatusCode) {
		return nil, fmt.Errorf("invalid HTTP request %d %s", resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}

func (d *dockerClient) DecodeStream(stream io.Reader) []string {
	type msg struct {
		Stream string `json:stream`
	}
	var msgs []string
	dec := json.NewDecoder(stream)
	for {
		var m msg
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		msgs = append(msgs, strings.TrimSuffix(m.Stream, "\n"))
	}
	return msgs
}

func (d *dockerClient) RemoveImage(name string, force bool, noprune bool) (io.ReadCloser, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("/images/%s", name)
		v      = &url.Values{}
	)
	if force {
		v.Set("force", "1")
	}
	if noprune {
		v.Set("noprune", "1")
	}
	uri = fmt.Sprintf("%s?%s", uri, v.Encode())

	respBody, err := d.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (d *dockerClient) ContainerWait(name string) error {
	var (
		method = "POST"
		uri    = fmt.Sprintf("/containers/%s/wait", name)
	)

	respBody, err := d.newRequest(method, uri, nil)
	if err != nil {
		return err
	}

	defer respBody.Close()

	ioutil.ReadAll(respBody)
	return nil
}

func (d *dockerClient) Attach(name string, logs, stream, stdout, stderr bool, inStream io.Writer) (io.Reader, io.Writer, error) {
	var (
		//method = "POST"
		uri = fmt.Sprintf("/containers/%s/attach", name)
	)
	var v url.Values
	if logs {
		v.Set("logs", "1")
	}
	if stream {
		v.Set("stream", "1")
	}
	if inStream != nil {
		v.Set("stdin", "1")
	}
	if stdout {
		v.Set("stdout", "1")
	}
	if stderr {
		v.Set("stderr", "1")
	}
	uri = fmt.Sprintf("%s?%s", uri, v.Encode())

	//respBody, conn, err := d.newRequest(method, uri, inStream)
	return nil, nil, nil
}

func (d *dockerClient) ContainerStats(name string) (io.ReadCloser, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("/containers/%s/stats", name)
	)

	respBody, err := d.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
