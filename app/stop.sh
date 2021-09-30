#!/bin/bash

if [ $# -ne 1 ] 
then 
    echo "Usage: ${0} ALGO"
    exit 1
fi

if [ $1 -eq 1 ] || [ $1 -eq 2 ] || [ $1 -eq 3 ]
then
    cd app_${1} && docker-compose -f algorithm-${1}.yml down && ./../clean_all.sh ${1}
else
    echo -e "The algorithm must be 1, 2 or 3."
fi
