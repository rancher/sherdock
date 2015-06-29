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

UI at http://localhost:8080

## Warning

Sherdock is a Work in Progress and running sherdock might lead to docker images being deleted on the host. The default 
configuration will not GC anything.  Please change the default configuration from ".*" to just the images you want to save.

## Developing

```bash

# Update UI
./script/build-ui

# Run
./script/run
```

## Release

    ./script/package
