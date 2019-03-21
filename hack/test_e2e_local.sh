#!/bin/bash

docker run -d -p 27017:27017  mongo

sleep 10

mongook=$(mongo --quiet --eval "db.runCommand({ping: 1})" | jq '.ok' 2> /dev/null )

if [ "$mongook" -eq "1" ]; then
  echo "Mongo looks OK"
else
  echo "Please start Mongo,,,"
  exit -1
fi


make clean
make build

ROOTDIR=$(git rev-parse --show-toplevel)


mongoimport -d miniplanes -c airports --type csv --file test/e2e/data/airports.dat --fieldFile=data/airports_schema.dat
mongoimport -d miniplanes -c airlines --type csv --file test/e2e/data/airlines.dat --fieldFile=data/airlines_schema.dat
mongoimport -d miniplanes -c courses --type csv --file test/e2e/data/courses.dat --fieldFile=data/courses_schema.dat

pushd schedules-generator/cmd/
go run main.go
popd

make clean
make build



#storage
${ROOTDIR}/_output/storage --mongo-host 127.0.0.1 --port 9999 &
STOPID=$!

#itineraries-server
${ROOTDIR}/_output/itineraries-server --storage-host  127.0.0.1 --storage-port 9999 --port 8888 &
ISPID=$!

${ROOTDIR}/_output/ui --itineraries-server-host 127.0.0.1 --itineraries-server-port 8888 --port 8080 --storage-host 127.0.0.1 --storage-port 9999 &
UIPID=$!

sleep 3

echo "Starting test"
cd ${ROOTDIR}/test/e2e && go test -c . && ./e2e.test --storage-host 127.0.0.1 --storage-port 9999 --itineraries-server-host 127.0.0.1 --itineraries-server-port 8888

kill $UIPID
kill $ISPID
kill $STOPID

docker rm $(docker stop $(docker ps -a -q --filter ancestor=mongo --format="{{.ID}}"))
