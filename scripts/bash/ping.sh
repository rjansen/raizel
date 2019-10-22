#!/usr/bin/env bash
name=$1
host=$2
port=$3
tries=$4
status="error"
for try in $(seq 0 $tries)
do
    nc -z $host $port && status="ok" && break
    printf "[%03d/%03d] ping %s on %s:%d\r" $try $tries $name $host $port
    sleep 1
done
printf "%s{%s:%d} %-25s\n" $name $host $port $status
if [[ "$status" == "error" ]]
then
    exit 5
fi
