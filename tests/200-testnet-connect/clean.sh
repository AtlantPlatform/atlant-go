#!/usr/bin/env bash
cd "$(dirname "$0")"
sudo docker-compose stop 
sudo docker-compose rm -f 
