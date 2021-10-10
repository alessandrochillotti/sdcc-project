#!/bin/bash

docker network rm app_my_net

docker rm app_node_1
docker rm app_node_2
docker rm app_node_3
docker rm app_sequencer_node_1
docker rm app_register_node_1
docker rm app_redis_node_1

docker rmi app_node
docker rmi app_sequencer_node
docker rmi app_register_node
