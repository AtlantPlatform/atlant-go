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
echo "[`date`] Starting Auth Server"
sudo docker-compose up -d --build auth
echo "[`date`] Starting Node 1"
sudo docker-compose up -d --build node1
sleep 4
NODE1=`curl -q http://localhost:33001/api/v1/ping || ''`
echo "[`date`] NODE1=$NODE1"
# TODO: give node permissions to write

echo "[`date`] Starting Node 2"
sudo docker-compose up -d --build node2
sleep 5
NODE1=`curl -q http://localhost:33002/api/v1/ping || ''`
echo "[`date`] NODE2=$NODE2"

echo "[`date`] Starting Node 3"
sudo docker-compose up -d --build node3
sleep 5
NODE1=`curl -q http://localhost:33003/api/v1/ping || ''`
echo "[`date`] NODE3=$NODE3"

echo "[`date`] Waiting for 10 seconds"
sleep 10

sudo docker-compose logs
sudo docker-compose down
sudo docker network rm clusterof3 || true