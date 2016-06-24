#!/bin/bash

set -e

go get github.com/constabulary/gb/...

export GOPATH=`pwd`:`pwd`/vendor
gb test