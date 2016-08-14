#!/bin/bash
#
# Copyright (c) 2016 Brian McKean
#
# Installs lamport in cloud
#
# Pre-req: must set up key pair for use with cloud service
#
set -x

# start nodes on cloud
#terraform apply -var num_nodes="3"
#sleep 120

# get host names from output of terraform script


fn="lamport_hosts.txt"

i="0"

while read -r line 
do 
    name="$line"
    #echo "$name"
    #echo "$i"
    arr[$i]="$name"
    i=$[$i+1]
done < "$fn"

num_servers="$i"

i="0"
while [ $i -lt $num_servers ]
do
    names[$i]=$(printf '%s\n' "$arr[[$i]" |  awk '{print $1}')
    #echo "${names[$i]}"
    ssh -o "StrictHostKeyChecking no" ubuntu@${names[$i]}  "/bin/bash" < lamport-setup.sh 1>${names[$i]}.out 2>${names[$i]}.err &
    i=$[$i+1]
done

wait

while [ $i -lt $num_servers ]
do
    ssh -o "StrictHostKeyChecking no" ubuntu@${names[$i]}  "/bin/bash" < lamport-node-init.sh &
    i=$[$i+1]
done

wait

echo "lamport-start done!"

