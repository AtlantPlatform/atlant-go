#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"  && pwd

cleanup() {
  docker-compose logs
  echo "NODE1=$NODE1"
  echo "NODE2=$NODE2"
  echo "NODE3=$NODE3"
  docker-compose stop
  docker-compose rm -f
  docker network rm clusterof3 || true  
}

atlant_client() {
  local node=$1
  docker-compose exec $node ./atlant-lite -A 127.0.0.1:33780 $2 $3 $4 $5 $6 $7 $8 $9
}

show_files() {
  local node=$1
  echo "[`date`] ======="
  echo "[`date`] Current Files Requesting /data at $node"
  atlant_client $node ls /data/  || true
  echo "[`date`] Displaying file /data/lipsum.txt"
  atlant_client $node get /data/lipsum.txt
  echo "[`date`] ======="
  echo 
}

show_stats() {
  local port=$1
  curl -s http://localhost:$port/api/v1/stats 
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
  docker-compose up -d --build $1
}

get_id() { 
    # port, resultvar
    local port=$1
    local __resultvar=$2
    local url="http://localhost:$port/api/v1/ping"
    echo "[`date`] Requesting URL $url"
    
    RESULT=`curl -s $url`
    echo "[`date`] Received $2: $RESULT"
    eval $__resultvar="$RESULT"
    expected $RESULT
}

wait_for_node() {
  local port=$1
  local name=$2
 
  let COUNTER=10
  # before continuing, ensure NODE0 exists
  while [[ "$(curl -s -o /dev/null -w '%{http_code}' http://localhost:$port/api/v1/ping)" != "200" ]]; do 
    echo "[`date`] Waiting for $name at localhost:$port... ($COUNTER)"
    COUNTER=$((COUNTER - 1))
    if [ "$COUNTER" == "0" ]; then
      echo "[`date`] TIMEOUT. $name NOT STARTED"
      cleanup
      exit 1
    fi
    sleep 5
  done
  get_id $port $name
}

wait_for_file() {
  local node=$1
  local fn=$2
  # 20 second timeout
  let COUNTER=6
  while [ "$RES_CLIENT" == "" ]; do 
    echo "[`date`] Waiting for file on the $node ($COUNTER)"
    echo "[`date`] Executing 'docker-compose exec $node ./atlant-lite -A 127.0.0.1:33780 get $fn'"
    RES_CLIENT=$(docker-compose exec $node ./atlant-lite -A 127.0.0.1:33780 get $fn) || true
    echo "[`date`] Result $RES_CLIENT"
    if [[ $RES_CLIENT == *"404 Not Found"* ]]; then
      RES_CLIENT=
    fi
    COUNTER=$((COUNTER - 1))
    if [ "$COUNTER" == "0" ]; then
      cleanup
      echo "[`date`] TIMEOUT ERROR. $node HAS NOT RECEIVED UPLOADED FILE"
      exit 1
    fi
    sleep 5
  done
}

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

# create network if not exists
docker network create --driver bridge clusterof3 || true

# starting the server
docker-compose up --build -d

wait_for_node 33001 "NODE1"
wait_for_node 33002 "NODE2"
wait_for_node 33003 "NODE3"

echo "[`date`] Uploading file to node1"
docker-compose exec node1 ./atlant-lite -A 127.0.0.1:33780 put ./lipsum.txt /data/lipsum.txt
echo "[`date`] Checking uploaded file on the 1st node"
docker-compose exec node1 ./atlant-lite -A 127.0.0.1:33780 get /data/lipsum.txt

wait_for_file node2 /data/lipsum.txt
wait_for_file node3 /data/lipsum.txt

# docker-compose exec -T node1 ./atlant-lite --help
echo "[`date`] Waiting for 10 seconds, collecting logs"
sleep 10

cleanup
echo "[`date`] Cluster of 3 nodes was successfully verified"