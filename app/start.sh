#!/bin/bash

# if less than two arguments supplied, display usage 
if [ $# -ne 1 ] 
then 
    echo "Usage: ${0} ALGO"
    exit 1
fi

if [ $1 -eq 1 ]
then
	echo "NODE_DOCKERFILE=dockerfile-1" > .env
	echo "REGISTER_DOCKERFILE=dockerfile-1" >> .env
    docker-compose --profile sequencer up --scale node=3
elif [ $1 -eq 2 ]
then
	echo "NODE_DOCKERFILE=dockerfile-2" > .env
	echo "REGISTER_DOCKERFILE=dockerfile-23" >> .env
	docker-compose --profile no_sequencer up --scale node=3
elif [ $1 -eq 3 ]
then
	echo "NODE_DOCKERFILE=dockerfile-3" > .env
	echo "REGISTER_DOCKERFILE=dockerfile-23" >> .env
	docker-compose --profile no_sequencer up --scale node=3
else
    echo -e "The algorithm must be 1, 2 or 3."
fi
