#!/bin/bash

source config.sh

scope=$1
coll=$2
field=$3
val=$4
query_port=$(($COUCHBASE_PORT+2))

index_query="create primary index $coll_idx on \`$COUCHBASE_BUCKET\`.\`$scope\`.\`$coll\`;"

curl -u $COUCHBASE_USER:$COUCHBASE_PASSWORD \
    -d "statement=$index_query" \
    http://localhost:$query_port/query/service

query="select dd from \`$COUCHBASE_BUCKET\`.\`_default\`.\`_default\` dd where $field=\"$val\";"

curl -u $COUCHBASE_USER:$COUCHBASE_PASSWORD \
    -d "statement=$query" \
    http://localhost:$query_port/query/service | jq --raw-output '.results[0:]' > res.json


cat res.json | jq '[.[] | (. * .dd) | del(.dd)]' > final_res.json
rm res.json
