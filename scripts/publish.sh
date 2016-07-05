#!/bin/bash

set -e

[ -z "$AUTOSCALE_BUCKET" ] && AUTOSCALE_BUCKET="mys3/autoscale"

gb build
mc cp bin/autoscalectl $AUTOSCALE_BUCKET

