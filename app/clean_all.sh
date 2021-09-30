#!/bin/bash

docker network rm app_1_my_net

docker rm app_1_node_1
docker rm app_1_node_2
docker rm app_1_node_3
docker rm app_1_sequencer_node_1
docker rm app_1_register_node_1

docker rmi app_1_node
docker rmi app_1_sequencer_node
docker rmi app_1_register_node
