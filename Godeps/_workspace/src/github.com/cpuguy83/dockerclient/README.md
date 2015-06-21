# Dockerclient

Dockerclient is a client library for Docker written in go, with a straight-forward
syntax.

### Usage

```go
package main

import (
  docker "github.com/cpuguy83/dockerclient"
  
  "fmt"
  "os"
  "strings"
)

func main() {
  client, err := docker.NewClient("tcp://127.0.0.1:2375")

  containers, err := client.FetchAllContainers(true)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  for _, container := range containers {
    container, err = client.FetchContainer(container.Id)
    if err != nil {
      fmt.Println(err)
    }

    name := strings.TrimPrefix(container.Name, "/")
    fmt.Println(name)
  }
}
```
