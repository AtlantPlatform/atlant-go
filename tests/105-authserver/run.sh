#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"  && pwd

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

cleanup() {
  docker-compose logs
  docker-compose stop && docker-compose rm -f 
}

# starting the server
docker-compose up -d --build
retVal=$?
if [ $retVal -ne 0 ]; then
    echo
    echo "[`date`] Build from Source Code Failed"
    exit $retVal
fi

# testing positive case:
# That node is up in 20 secs
let COUNTER=4
# before continuing, ensure NODE0 exists
while [[ "$(curl -s -o /dev/null -w '%{http_code}' http://localhost:33700/)" != "200" ]]; do 
  echo "[`date`] Waiting for auth server... ($COUNTER)"
  COUNTER=$((COUNTER - 1))
  if [ "$COUNTER" == "0" ]; then
    cleanup
    echo "[`date`] TIMEOUT ERROR. auth node NOT STARTED"
    exit 1
  fi
  sleep 5
done

# That node can write and result is the same
EXPECTED="14V21784834838489:write,sync"
curl -q http://localhost:33700/ -d "$EXPECTED"

RESULT=`curl -sq -X GET http://localhost:33700/`
echo "RESULT1:'$RESULT'"
if [ "$RESULT" != "$EXPECTED" ]; then
    echo "EXPECTED:'$EXPECTED'"
    exit 1
fi
echo "RESPONSE OK"
EXPECT_FORBIDDEN=`curl -sq -H"X-Forwarded-For: 8.8.8.8" http://localhost:33700/ -d "$EXPECTED"`
if [ "$EXPECT_FORBIDDEN" != "No permission to update" ]; then
    echo "ERROR: Expected to be forbidden because of external address"
    exit 1
fi

cleanup
echo "[`date`] AUTH NODE was successfully verified"