#!/bin/bash


mongook=$(mongo --quiet --eval "db.runCommand({ping: 1})" | jq '.ok' 2> /dev/null )

if [ "$mongook" -eq "1" ]; then
  echo "Mongo looks OK"
else
  echo "Please start Mongo,,,"
  exit -1
fi



ROOTDIR=$(git rev-parse --show-toplevel)

#TODO check mongo is running

#e2e storage
${ROOTDIR}/_output/storage --mongo-host 127.0.0.1 --port 9999 &
REGPFPID=$!

sleep 5

cd ${ROOTDIR}/test/e2e && go test -c . && ./e2e.test --storage-host 127.0.0.1 --storage-port 9999

kill $REGPFPID
