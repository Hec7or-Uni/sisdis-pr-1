#!/bin/bash

# Ping a list of hosts and save the ones that are up to a file
# Usage: ping.sh <file with list of hosts> <desire port> <numer of hosts>

LABS=$1
PORT=$2
MAX_HOSTS=$3
NIP=$4


# Check if the file exists
if [ ! -f ${LABS} ]; then
    echo "File ${LABS} does not exist"
    exit 1
else 
    rm "endpoints.txt"
fi

# Check if the file is empty
if [ ! -s ${LABS} ]; then
    echo "File ${LABS} is empty"
    exit 1
fi

# Check if the file is readable
if [ ! -r ${LABS} ]; then
    echo "File ${LABS} is not readable"
    exit 1
fi

# Check if the output file exists
if [ -f $2 ]; then
    rm $2
fi

# Loop through the list of hosts until reach the max number of hosts
num_hosts=0
while read _ host; do
    echo "Pinging $host... "
    
    # Check max number of hosts
    if [ $((${MAX_HOSTS} == $num_hosts)) -eq 1 ]; then
        echo "Max number of hosts reached."
        break
    fi

    ping -c 1 $host > /dev/null 2>&1
    # check if the host is up
    if [ $? -eq 0 ]; then
        echo "$host is up"
        num_hosts=$((num_hosts+1))
        echo ${host}:$((${PORT}+num_hosts)) >> "endpoints.txt"
    else
        echo "$host is down"
    fi
done < ${LABS}

# run all workers
IFS=:
while  read ip port; do
    echo $endpoint
    # execute remote command
    echo ssh -i "~/.ssh/id_rsa" a${NIP}@${ip} "go run worker.go ${ip} ${port}"
done < "endpoints.txt"