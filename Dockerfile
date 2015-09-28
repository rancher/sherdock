FROM golang:1.4
EXPOSE 8080
RUN go get github.com/tools/godep
COPY . /go/src/github.com/rancher/sherdock
WORKDIR /go/src/github.com/rancher/sherdock
RUN godep go build -ldflags "-linkmode external -extldflags -static" -o /usr/bin/sherdock /go/src/github.com/rancher/sherdock/main.go /go/src/github.com/rancher/sherdock/bindata_assetfs.go
CMD ["sherdock"]
