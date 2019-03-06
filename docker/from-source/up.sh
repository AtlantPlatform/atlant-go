#!/bin/bash
if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

# starting the server
sudo docker-compose up -d --build
sudo docker-compose logs -f