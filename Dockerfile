FROM golang:1.4
COPY main.go /main.go
RUN go get github.com/cpuguy83/dockerclient
RUN go get github.com/emicklei/go-restful
RUN go get github.com/samalba/dockerclient
CMD go run /main.go
EXPOSE 8080
