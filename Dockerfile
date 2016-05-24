FROM alpine:edge
MAINTAINER Bryan Liles <bliles@digitalocean.com>

ENV PATH /go/bin:$PATH
ENV GOPATH /go:/go/vendor

RUN apk add --update-cache bash go ca-certificates

ADD . /go
WORKDIR /go
RUN CGO_ENABLED=0 go install autoscale/cmd/do-autoscale
RUN apk del -v go && \
  rm -rf /var/cache/apk/*

ENTRYPOINT /go/bin/do-autoscale
