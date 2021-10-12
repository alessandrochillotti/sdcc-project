#!/bin/bash

docker network rm app_my_net

docker rm app_node_1
docker rm app_node_2
docker rm app_node_3
docker rm app_node_4
docker rm app_node_5
docker rm app_node_6
docker rm app_node_7
docker rm app_node_8
docker rm app_node_9
docker rm app_node_10
docker rm app_sequencer_node_1
docker rm app_register_node_1

docker rmi app_node
docker rmi app_sequencer_node
docker rmi app_register_node
