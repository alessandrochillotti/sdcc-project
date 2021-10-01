#!/bin/bash

docker network rm app_${1}_my_net

docker rm app_${1}_node_1
docker rm app_${1}_node_2
docker rm app_${1}_node_3
docker rm app_${1}_sequencer_node_1
docker rm app_${1}_register_node_1
docker rm app_2_redis_node_1

docker rmi app_${1}_node
docker rmi app_${1}_sequencer_node
docker rmi app_${1}_register_node
