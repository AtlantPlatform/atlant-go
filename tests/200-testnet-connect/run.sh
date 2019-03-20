#!/usr/bin/env bash
set -e
cd "$(dirname "$0")" && pwd

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

# starting the server
sudo docker-compose up -d --build
retVal=$?
if [ $retVal -ne 0 ]; then
    echo
    echo "[`date`] Build from Source Code Failed"
    exit $retVal
fi

echo "[`date`] Waiting for 7 seconds"
sleep 7
echo "[`date`] Checking ping"
PING_RESULT=$(curl -v -s http://localhost:33780/api/v1/ping | tail -1)
echo "[`date`] Received ping: $PING_RESULT"

sudo docker-compose logs
sudo docker-compose down

if [ -z "$PING_RESULT" ]; then
    echo "[`date`] Start Verification Failed"
    exit 1
fi

echo "[`date`] Build and start of 'atlant-go' was successfully verified, congrats"