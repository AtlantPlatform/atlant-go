#!/usr/bin/env bash
set -e
cd "$(dirname "$0")" && pwd

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

cleanup() {
  sudo docker-compose logs
  sudo docker-compose stop && sudo docker-compose rm -f
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

# 30 second timeout
let COUNTER=10
# before continuing, ensure TESTNET NODE exists
while [[ "$(curl -s -o /dev/null -w '%{http_code}' http://localhost:33780/api/v1/ping)" != "200" ]]; do 
  echo "[`date`] Waiting for testnet node... ($COUNTER)"
  COUNTER=$((COUNTER - 1))
  if [ "$COUNTER" == "0" ]; then
    cleanup
    echo "[`date`] TIMEOUT ERROR. TEST NODE NOT STARTED"
    exit 1
  fi
  sleep 5
done

get_value "ping"
get_value "newID"
get_value "env"
get_value "session"
get_value "version"

echo "[`date`] Getting index"
curl -s http://localhost:33780/index/

if [ "$1" != "stay" ]; then
  cleanup
  echo "[`date`] Build and start of 'atlant-go' was successfully verified, congrats"
else
  echo "[`date`] Build and start of 'atlant-go' was successfully verified, node is still running"
fi