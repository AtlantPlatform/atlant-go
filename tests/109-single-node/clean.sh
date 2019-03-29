#!/usr/bin/env bash
cd "$(dirname "$0")"
sudo docker-compose stop 
sudo docker-compose rm -f 
sudo docker network rm clusterof1 || true