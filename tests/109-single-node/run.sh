#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"  && echo "[`date`] `pwd`"

cleanup() {
  sudo docker-compose logs
  echo "NODE0=$NODE0"
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


if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

# create network if not exists
sudo docker network create --driver bridge clusterof1 || true

# starting the server
sudo docker-compose up --build -d
sleep 5

# before continuing, ensure NODE1 exists
# get_id 33101 NODE0
while [[ "$(curl -s -o /dev/null -w '%{http_code}' http://localhost:33101/api/v1/ping)" != "200" ]]; do 
  echo "[`date`] Waiting for node0..."
  sleep 5; 
done

echo "[`date`] Uploading file"
sudo docker-compose exec node0 ./atlant-lite -A 127.0.0.1:33780 put ./lipsum.txt /data/lipsum.txt
echo "[`date`] Checking uploaded file"
sudo docker-compose exec node0 ./atlant-lite -A 127.0.0.1:33780 get /data/lipsum.txt

# sudo docker-compose exec -T node1 ./atlant-lite --help
echo "[`date`] Waiting for 10 seconds, collecting logs"
sleep 10

cleanup
echo "[`date`] Cluster of 3 nodes was successfully verified"