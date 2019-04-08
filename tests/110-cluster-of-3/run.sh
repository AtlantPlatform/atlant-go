#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"  && pwd

cleanup() {
  sudo docker-compose logs
  echo "NODE1=$NODE1"
  echo "NODE2=$NODE2"
  echo "NODE3=$NODE3"
  sudo docker-compose stop
  sudo docker-compose rm -f
  sudo docker network rm clusterof3 || true  
}

atlant_client() {
  local node=$1
  sudo docker-compose exec $node ./atlant-lite -A 127.0.0.1:33780 $2 $3 $4 $5 $6 $7 $8 $9
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
sudo docker-compose up --build -d
sleep 30

# before continuing, ensure NODE1 exists
while [[ "$(curl -s -o /dev/null -w '%{http_code}' http://localhost:33001/api/v1/ping)" != "200" ]]; do 
  echo "[`date`] Waiting for node1..."
  sleep 5; 
done
get_id 33001 NODE1


# before continuing, ensure NODE2 exists
while [[ "$(curl -s -o /dev/null -w '%{http_code}' http://localhost:33002/api/v1/ping)" != "200" ]]; do 
  echo "[`date`] Waiting for node2..."
  sleep 5; 
done
get_id 33002 NODE2

# before continuing, ensure NODE3 exists
while [[ "$(curl -s -o /dev/null -w '%{http_code}' http://localhost:33003/api/v1/ping)" != "200" ]]; do 
  echo "[`date`] Waiting for node3..."
  sleep 5; 
done
get_id 33003 NODE3

echo "[`date`] Uploading file to node1"
sudo docker-compose exec node1 ./atlant-lite -A 127.0.0.1:33780 put ./lipsum.txt /data/lipsum.txt
echo "[`date`] Checking uploaded file on the 1st node"
sudo docker-compose exec node1 ./atlant-lite -A 127.0.0.1:33780 get /data/lipsum.txt

echo "[`date`] Waiting for 5 seconds"
sleep 5
echo "[`date`] Checking uploaded file on the 2nd node"
sudo docker-compose exec node2 ./atlant-lite -A 127.0.0.1:33780 get /data/lipsum.txt
echo "[`date`] Checking uploaded file on the 3rd node"
sudo docker-compose exec node3 ./atlant-lite -A 127.0.0.1:33780 get /data/lipsum.txt

# sudo docker-compose exec -T node1 ./atlant-lite --help
echo "[`date`] Waiting for 10 seconds, collecting logs"
sleep 10

cleanup
echo "[`date`] Cluster of 3 nodes was successfully verified"