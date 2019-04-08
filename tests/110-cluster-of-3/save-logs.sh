#!/usr/bin/env bash
cd "$(dirname "$0")"
sudo docker-compose logs node1 2>&1 | sed 's/\x1b\[[0-9;]*m//g' > node1.log
sudo docker-compose logs node2 2>&1 | sed 's/\x1b\[[0-9;]*m//g' > node2.log
sudo docker-compose logs node3 2>&1 | sed 's/\x1b\[[0-9;]*m//g' > node3.log
