#!/bin/bash

# make sure the sqlite data file is there
if [ ! -e "/home/ubuntu/data/tithedeclare.db" ]; then mkdir -p "/home/ubuntu/data"; touch "/home/ubuntu/data/tithedeclare.db"; fi

# build and run docker
docker build -t tithedeclare:latest -f ./build/Dockerfile .
is_running=$(docker ps | grep tithedeclare | wc -l | xargs)
is_container=$(docker ps -a | grep tithedeclare | wc -l | xargs)

if [ "$is_running" == "1" ]; then docker stop tithedeclare; fi 
if [ "$is_container" == "1" ]; then docker rm tithedeclare; fi

docker run -d --name=tithedeclare -v /home/ubuntu/data:/app/data -p 12672:12580 --env-file ./env_vars tithedeclare:latest