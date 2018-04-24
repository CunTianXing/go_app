#!/bin/bash

function dkcl(){
        CONTAINER_IDS=$(docker ps -aq)
	echo
        if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" = " " ]; then
                echo "========== No containers available for deletion =========="
        else
                echo "CONTAINER_IDS=========>>>"
                echo $CONTAINER_IDS
                docker rm -f $CONTAINER_IDS
        fi
	echo
}

function dkrm(){
        DOCKER_IMAGE_IDS=$(docker images | grep "fabsdkgo\|none\|v[0-9]-" | awk '{print $3}')
	echo
        if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" = " " ]; then
		echo "========== No images available for deletion ==========="
        else
                echo "DOCKER_IMAGE_IDS=========>>>"
                echo $DOCKER_IMAGE_IDS
                docker rmi -f $DOCKER_IMAGE_IDS
        fi
	echo
}

function restartNetwork(){
  echo "restart"
  docker-compose down
  dkcl
  dkrm
  rm -rf /tmp/msp/keystore/
  #Start the network
  docker-compose up -d
  echo "ok"
}
if [ "$1" == "up" ]; then
	docker-compose up -d
elif [ "$1" == "down" ]; then ## Clear the network
docker-compose down
dkcl
dkrm
rm -rf /tmp/msp/keystore/
elif [ "$1" == "restart" ]; then ## Restart the network
	restartNetwork
else
  echo "no params"
	exit 1
fi
#restartNetwork
