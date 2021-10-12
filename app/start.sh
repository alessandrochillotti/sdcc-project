#!/bin/bash

# if less than two arguments supplied, display usage 
if [ $# -ne 2 ] 
then 
    echo "Usage: ${0} ALGO NODES"
    exit 1
fi

if [ $1 -eq 1 ]
then
    ALGO=1 NUMBER_NODE=${2} docker-compose --profile sequencer up --scale node=${2} -d
elif [ $1 -eq 2 ]
then
	ALGO=2 NUMBER_NODE=${2} docker-compose --profile no_sequencer up --scale node=${2} -d
elif [ $1 -eq 3 ]
then
	ALGO=3 NUMBER_NODE=${2} docker-compose --profile no_sequencer up --scale node=${2} -d
else
    echo -e "The algorithm must be 1, 2 or 3."
fi
