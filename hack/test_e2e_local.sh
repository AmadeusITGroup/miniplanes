#!/bin/bash


mongook=$(mongo --quiet --eval "db.runCommand({ping: 1})" | jq '.ok' 2> /dev/null )

if [ "$mongook" -eq "1" ]; then
  echo "Mongo looks OK"
else
  echo "Please start Mongo,,,"
  exit -1
fi

ROOTDIR=$(git rev-parse --show-toplevel)

#storage
${ROOTDIR}/_output/storage --mongo-host 127.0.0.1 --port 9999 &
STOPID=$!

#itineraries-server
${ROOTDIR}/_output/itineraries-server --storage-host  127.0.0.1 --storage-port 9999 --port 8888 &
ISPID=$!

sleep 3

echo "Starting test"
cd ${ROOTDIR}/test/e2e && go test -c . && ./e2e.test --storage-host 127.0.0.1 --storage-port 9999 --itineraries-server-host 127.0.0.1 --itineraries-server-port 8888

kill $ISPID
kill $STOPID
