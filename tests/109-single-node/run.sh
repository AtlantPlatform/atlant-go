#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"  && echo "[`date`] `pwd`"

cleanup() {
  sudo docker-compose logs
  sudo docker-compose stop && sudo docker-compose rm -f && sudo docker network rm clusterof1 || true
}

expected() {
  if [ "$1" == "" ]; then
    cleanup
    echo "[CLUSTER TESTS FAILED]"
    exit 1
  fi
}

start_service() {
  echo "[`date`] Starting $1"
  sudo docker-compose up -d --build $1
}


get_id() { 
    # port, resultvar
    local port=$1
    local __resultvar=$2
    local url="http://localhost:$port/api/v1/ping"
    echo "[`date`] Requesting URL $url"
    
    RESULT=`curl $url`
    echo "[`date`] Received $2: $RESULT"
    eval $__resultvar="$RESULT"
    expected $RESULT
}

atlant_client() {
    sudo docker-compose exec node0 ./atlant-lite -A 127.0.0.1:33780 $@
}

show_files() {
    echo "[`date`] ======="
    echo "[`date`] Current Files Requesting /data"
    atlant_client ls /data/  || true
    echo "[`date`] Displaying file /data/lipsum.txt"
    atlant_client get /data/lipsum.txt
    echo "[`date`] ======="
    echo 
}

show_stats() {
    curl -s http://localhost:33101/api/v1/stats 
}

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

# create network if not exists
sudo docker network create --driver bridge clusterof1 || true

# starting the server
sudo docker-compose up --build -d
sleep 5

# 30 second timeout
let COUNTER=6
# before continuing, ensure NODE0 exists
while [[ "$(curl -s -o /dev/null -w '%{http_code}' http://localhost:33101/api/v1/ping)" != "200" ]]; do 
  echo "[`date`] Waiting for node0... ($COUNTER)"
  COUNTER=$((COUNTER - 1))
  if [ "$COUNTER" == "0" ]; then
    cleanup
    echo "[`date`] TIMEOUT ERROR. NODE0 NOT STARTED"
    exit 1
  fi
  sleep 5
done

echo "[`date`] Uploading file"

atlant_client put ./lipsum.txt /data/lipsum.txt
show_files
atlant_client put ./lipsum_v2.txt /data/lipsum.txt
show_files
atlant_client versions /data/lipsum.txt

echo 
export RECORDID=`sudo docker-compose exec node0 ./atlant-lite -A 127.0.0.1:33780 meta /data/lipsum.txt 2>&1 | grep '"id"' | cut -f4 -d'"'`
echo "[`date`] RECORDID='$RECORDID'"
echo 
echo "[`date`] Deleting record"
atlant_client delete $RECORDID

# echo "[`date`] Deleting successful"
# show_files
# echo "[`date`] Waiting for 10 seconds, collecting logs"
# sleep 10

cleanup
echo "[`date`] CREATE+GET case was successfully verified with single node"