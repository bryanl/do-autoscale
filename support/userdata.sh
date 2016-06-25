#!/bin/bash

ARCHIVE_BASE=https://s3.pifft.com/autoscale

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

apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" > /etc/apt/sources.list.d/docker.list
apt-get update
apt-get install -q -y docker-engine

docker pull bryanl/do-autoscale
docker pull postgres
docker pull prom/prometheus

for i in autoscale postgres prometheus; do
  echo "downloading ${i}..."
  curl -sSL -o /tmp/$i.service ${ARCHIVE_BASE}/$i.service
  mv /tmp/$i.service /etc/systemd/system
done

systemctl daemon-reload

for i in autoscale postgres prometheus; do
  echo "starting ${i}..."
  systemctl start $i
done