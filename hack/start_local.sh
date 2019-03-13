#!/bin/bash

mongook=$(mongo --quiet --eval "db.runCommand({ping: 1})" | jq '.ok' 2> /dev/null )

if [ "$mongook" -eq "1" ]; then
  echo "Mongo looks OK"
else
  echo "Please start Mongo,,,"
  exit -1
fi

#make clean
#make build

ROOTDIR=$(git rev-parse --show-toplevel)

#storage
${ROOTDIR}/_output/storage --mongo-host 127.0.0.1 --port 9999 &
STOPID=$!

#itineraries-server
${ROOTDIR}/_output/itineraries-server --storage-host  127.0.0.1 --storage-port 9999 --port 8888 &
ISPID=$!

${ROOTDIR}/_output/ui --itineraries-server-host 127.0.0.1 --itineraries-server-port 8888 --port 8080 --storage-host 127.0.0.1 --storage-port 9999 &
UIPID=$!

sleep 3

read  -n 1 -p "Enter to stop..." input


kill $UIPID
kill $ISPID
kill $STOPID
