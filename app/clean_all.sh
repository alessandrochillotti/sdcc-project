#!/bin/bash

docker network rm app_my_net

docker rm app_peer_1
docker rm app_peer_2
docker rm app_peer_3
docker rm app_peer_4
docker rm app_peer_5
docker rm app_peer_6
docker rm app_peer_7
docker rm app_peer_8
docker rm app_peer_9
docker rm app_peer_10
docker rm app_sequencer_node_1
docker rm app_register_node_1

docker rmi app_peer
docker rmi app_sequencer_node
docker rmi app_register_node
