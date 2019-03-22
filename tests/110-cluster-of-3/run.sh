#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"  && pwd

cleanup() {
  sudo docker-compose logs
  echo "NODE1=$NODE1"
  echo "NODE2=$NODE2"
  echo "NODE3=$NODE3"
  sudo docker-compose down
  sudo docker network rm clusterof3 || true  
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


if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

# create network if not exists
sudo docker network create --driver bridge clusterof3 || true

# starting the server
start_service auth
start_service node1
sleep 15
get_id 33001 NODE1

echo "[`date`] Giving Node 1 permission to write"
curl http://localhost:33700/ -d "$NODE1:write,sync"

start_service node2
start_service node3

sleep 30
get_id 33002 NODE2
get_id 33003 NODE3

# sudo docker-compose exec -T node1 ./atlant-lite --help
echo "[`date`] Waiting for 10 seconds, collecting logs"
sleep 10

cleanup
echo "[`date`] Cluster of 3 nodes was successfully verified"