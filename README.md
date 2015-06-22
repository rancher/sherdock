Sherdock
========

![sherdock](logo.png "SherDock")

DockerCon 2015 Hackathon Project

## Features

* Automatic GC of images based on regexp
* Find and delete orphan Docker volumes (requires Docker 1.7)
* UI

## Running

    docker run -d -v /var/lib/docker:/var/lib/docker -v /var/run/docker.sock:/var/run/docker.sock -p 8080:8080 rancher/sherdock

UI at http://localhost:8008

## Warning

Sherdock is a Work in Progress and running sherdock might lead to docker images being deleted on the host. The default 
configuration keeps only 3 images: `ubuntu:latest`, `busybox:latest` and `rancher/server`, so all other images will be deleted. Run at your own risk.

## Developing

```bash

# Update UI
./script/build-ui

# Run
./script/run
```

## Release

    ./script/package
