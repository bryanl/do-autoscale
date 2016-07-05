#!/bin/bash

ARCHIVE_BASE=https://s3.pifft.com/autoscale
CTL_BIN=/usr/local/bin/autoscalectl

error() {
  local parent_lineno="$1"
  local message="$2"
  local code="${3:-1}"
  if [[ -n "$message" ]] ; then
    echo "Error on or near line ${parent_lineno}: ${message}; exiting with status ${code}"
  else
    echo "Error on or near line ${parent_lineno}; exiting with status ${code}"
  fi
  exit "${code}"
}
trap 'error ${LINENO}' ERR

# download and install service files

curl -sSL -o ${CTL_BIN} ${ARCHIVE_BASE}/autoscalectl
chmod +x ${CTL_BIN}
