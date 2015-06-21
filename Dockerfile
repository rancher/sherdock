FROM golang:1.4
EXPOSE 8080
RUN go get github.com/cpuguy83/dockerclient \
  github.com/emicklei/go-restful \
  github.com/samalba/dockerclient \
  github.com/hashicorp/golang-lru
COPY . /go/src/github.com/rancherio/sherdock
CMD go run /go/src/github.com/rancherio/sherdock/main.go \
  /go/src/github.com/rancherio/sherdock/images.go \
  /go/src/github.com/rancherio/sherdock/containers.go
