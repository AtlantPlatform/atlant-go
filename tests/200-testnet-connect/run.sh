#!/usr/bin/env bash
set -e
cd "$(dirname "$0")" && pwd

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

cleanup() {
  sudo docker-compose logs
  sudo docker-compose down
}

expected() {
  if [ "$1" == "" ]; then
    cleanup
    echo "FAILURE..."
    exit 1
  fi
}

get_value() {
    export RESULT=`curl -s http://localhost:33780/api/v1/$1 || ''`
    echo "[`date`] Received $1: $RESULT"
    expected $RESULT
}

# starting the server
sudo docker-compose up -d --build
echo "[`date`] Waiting for 8 seconds"
sleep 8

get_value "ping"
get_value "newID"
get_value "env"
get_value "session"
get_value "version"

sudo docker-compose exec atlant_testnet ./atlant-lite -A 127.0.0.1:33780 debug

# cleanup
# echo "[`date`] Build and start of 'atlant-go' was successfully verified, congrats"