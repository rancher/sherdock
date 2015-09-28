FROM golang:1.4
EXPOSE 8080
RUN go get github.com/tools/godep
RUN apt-get update && apt-get install -y g++
RUN curl -L http://nodejs.org/dist/v0.12.4/node-v0.12.4-linux-x64.tar.gz | tar xvzf - -C /usr/local --strip-components=1
RUN npm install -g ember-cli
RUN npm install -g bower
RUN go get github.com/jteeuwen/go-bindata/...
RUN go get github.com/elazarl/go-bindata-assetfs/...
COPY sherdock-ember /sherdock-ember
RUN cd /sherdock-ember && \
    npm install
RUN cd /sherdock-ember && \
    bower install --allow-root
RUN cd /sherdock-ember && \
    yes | ember update && \
    ember build
COPY . /go/src/github.com/rancher/sherdock
WORKDIR /go/src/github.com/rancher/sherdock
RUN cp -rf /sherdock-ember/dist /go/src/github.com/rancher/sherdock/data && \
    cd /go/src/github.com/rancher/sherdock && \
    for i in data/images data/graph data/volumes data/config; do mkdir -p $i && cp data/index.html $i; done && \
    go-bindata-assetfs data/...
RUN godep go build -o /usr/bin/sherdock /go/src/github.com/rancher/sherdock/main.go /go/src/github.com/rancher/sherdock/bindata_assetfs.go
CMD sherdock
