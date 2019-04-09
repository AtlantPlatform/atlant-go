#!/usr/bin/env bash
cd "$(dirname "$0")"
docker-compose stop 
docker-compose rm -f 
docker network rm clusterof1 || true