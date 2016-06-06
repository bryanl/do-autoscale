FROM golang:1.6-alpine
MAINTAINER Bryan Liles <bliles@digitalocean.com>

ENV PATH /go/bin:$PATH
ENV GOPATH /go:/go/vendor

ADD . /go
WORKDIR /go
RUN CGO_ENABLED=0 go install autoscale/cmd/do-autoscale

ENTRYPOINT ["/go/bin/do-autoscale"]
