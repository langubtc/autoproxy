FROM golang:latest
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /gopath/
ENV GOPATH=/gopath/
ENV GOOS=linux

RUN go get -u -v github.com/lixiangyun/autoproxy
WORKDIR /gopath/src/github.com/lixiangyun/autoproxy/proxy
RUN go build .

FROM ubuntu:xenial
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /usr/bin/
COPY --from=0 /gopath/src/github.com/lixiangyun/autoproxy/proxy/proxy ./autoproxy
COPY --from=0 /gopath/src/github.com/lixiangyun/autoproxy/proxy/server.yaml ./config.yaml

RUN chmod +x autoproxy

EXPOSE 8080

ENTRYPOINT ["autoproxy"]
