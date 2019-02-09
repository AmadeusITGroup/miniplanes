#!/bin/bash


ROOTDIR=$(git rev-parse --show-toplevel)

#e2e storage
${ROOTDIR}/_output/storage --mongo-host 127.0.0.1 --port 9999 &
REGPFPID=$!

sleep 5

cd ${ROOTDIR}/test/e2e && go test -c . && ./e2e.test --storage-host 127.0.0.1 --storage-port 9999




kill $REGPFPID
