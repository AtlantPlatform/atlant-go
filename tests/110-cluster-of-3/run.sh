#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"  && pwd

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

# create network if not exists
sudo docker network create --driver bridge clusterof3 || true

# starting the server
sudo docker-compose up -d --build
retVal=$?
if [ $retVal -ne 0 ]; then
    echo
    echo "[`date`] Build from Source Code Failed"
    exit $retVal
fi

echo "[`date`] Waiting for 10 seconds"
sleep 10

sudo docker-compose logs
sudo docker-compose down
sudo docker network rm clusterof3 || true