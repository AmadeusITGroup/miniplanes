#!/bin/bash

#To be run like ./hack/post_json_data.sh data/airlines_schema.dat data/airlines.dat  http://127.0.0.1:9999/airlines

schema=$1
data=$2
URL=$3

SCHEMA=()
while read -r schema_line; do
    SCHEMA+=(\"${schema_line^}\")
done < "$schema"

while IFS=, read -ra data_line; do
  if [ ${#data_line[@]} -eq ${#SCHEMA[@]} ]; then #minimum check wether data fit schema
    for ((i=0;i<${#data_line[@]};++i)); do
      if [ "${data_line[i]}" = "\N" ]; then # filter dirty field in data
        data_line[i]=\"\"
      fi
    done

    CURLCMD="curl -s -o /dev/null -w \"%{http_code}\" --header \"Content-Type: application/json\" --request POST --data '"
    CURLCMD+="{"
    for ((i=0;i<${#SCHEMA[@]};++i)); do
      CURLCMD+="${SCHEMA[i]}: ${data_line[i]}"
      COUNT=$((${#SCHEMA[@]} - 1))
      if [ $i -lt ${COUNT} ]; then
        CURLCMD+=", "
      fi
   done
    CURLCMD+="}' $URL"
    httpcode=$(eval ${CURLCMD})
    if [ "$httpcode" != "201" ]; then
       echo "ERROR FOR ${CURLCMD} ${httpcode}"
    fi
  fi
done < "$data"
