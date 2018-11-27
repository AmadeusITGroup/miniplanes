#!/bin/bash

function contains() {
    local n=$#
    local value=${!n}
    for ((i=1;i < $#;i++)) {
        if [ "${!i}" == "${value}" ]; then
            echo "y"
            return 0
        fi
    }
    echo "n"
    return 1
}

#Fetch from local airlines dat file
KNOWNAIRLINES=($(OFS="," awk -F,  '{gsub(/"/, "", $4); printf("%s,%s ", $4, $1)}' airlines.dat))

tmproutes=$(mktemp /tmp/tmproutes.XXXXXX)
#Generate all temporary routes
for a in "${KNOWNAIRLINES[@]}"
do
   :
   grep ^$a ../../../data/routes.dat >> "$tmproutes"
done

KNOWNAIRPORTS=($(awk -F, '{print $1}' airports.dat ))

while read -r p; do
  source=$(echo $p | awk -F, '{print $4}')
  dest=$(echo $p | awk -F, '{print $6}')
  if [ $(contains "${KNOWNAIRPORTS[@]}" $source) == "y" -a $(contains "${KNOWNAIRPORTS[@]}" $dest) == "y"  ]; then
      echo $p >> routes.dat 
  fi
done < "$tmproutes"

rm "$tmproutes"