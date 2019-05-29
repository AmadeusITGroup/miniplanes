#!/bin/bash

#To be run like ./hack/submit_data.sh data/airlines_schema.dat data/airlines.dat  http://127.0.0.1:8080/save_airline
#curl -u admin:admin --dat "AirlineID=21317&Name=Svyaz%20Rossiya&Alias=Russian%20Commuter%20&IATA=7R&ICAO=SJM&Callsign=RussianConnecty&Country=Russia&Active=Y" http://127.0.0.1:8080/save_airline

schema=$1
data=$2
URL=$3

SCHEMA=()
while read -r schema_line; do
    SCHEMA+=(${schema_line^})
done < "$schema"

while IFS=, read -ra data_line; do
  if [ ${#data_line[@]} -eq ${#SCHEMA[@]} ]; then #minimum check wether data fit schema
    for ((i=0;i<${#data_line[@]};++i)); do
      #echo "original data_line ${data_line[i]}"
      temp=${data_line[i]%\"}
      #echo "first filter temp -> $temp"
      temp=${temp#\"}
      data_line[i]=$temp
      if [ "${data_line[i]}" = "\N" ]; then # filter dirty field in data
        data_line[i]=""
      fi
      done

    CURLCMD="curl -u admin:admin --data \""
    for ((i=0;i<${#SCHEMA[@]};++i)); do
      CURLCMD+=${SCHEMA[i]}=${data_line[i]}
      COUNT=$((${#SCHEMA[@]} - 1))
      if [ $i -lt ${COUNT} ]; then
        CURLCMD+="&"
      fi
     done
     CURLCMD+="\" $URL"
     echo $CURLCMD
     httpcode=$(eval ${CURLCMD})
     fi
done < "$data"
