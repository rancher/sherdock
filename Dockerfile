FROM golang:1.4
EXPOSE 8080
RUN go get github.com/cpuguy83/dockerclient \
  github.com/emicklei/go-restful \
  github.com/samalba/dockerclient \
  github.com/hashicorp/golang-lru
COPY . /go/src/github.com/rancherio/sherdock
WORKDIR /go/src/github.com/rancherio/sherdock
RUN go build -o /usr/bin/sherdock /go/src/github.com/rancherio/sherdock/main.go
CMD sherdock
