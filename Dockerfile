FROM alpine:edge
MAINTAINER Bryan Liles <bliles@digitalocean.com>

ENV PATH /go/bin:$PATH
ENV GOPATH /go:/go/vendor
ADD . /go
WORKDIR /go

RUN apk add --update-cache bash go bzr git mercurial subversion openssh-client ca-certificates && \
  CGO_ENABLED=0 go install autoscale/cmd/do-autoscale && \
  apk del -v go bzr git mercurial subversion openssh-client && \
  rm -rf /var/cache/apk/*

