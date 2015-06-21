FROM golang:1.4
EXPOSE 8080
RUN go get github.com/tools/godep
COPY . /go/src/github.com/rancherio/sherdock
WORKDIR /go/src/github.com/rancherio/sherdock
RUN godep go build -o /usr/bin/sherdock /go/src/github.com/rancherio/sherdock/main.go /go/src/github.com/rancherio/sherdock/bindata_assetfs.go
CMD sherdock
