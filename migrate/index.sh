#!/bin/bash

source config.sh

scope=$1
coll=$2
index=$3
shift
shift
shift
#arr is an array of the columns to be indexed
arr=("$@")

col_str=$(printf ",%s" "${arr[@]}")
col_str=${col_str:1}
echo $col_str

query_port=$(($COUCHBASE_PORT+2))

curl -u $COUCHBASE_USER:$COUCHBASE_PASSWORD \
    -d "statement=create index \`$index\` on \`$COUCHBASE_BUCKET\`.\`$scope\`.\`$coll\`(\`$col_str\`);" \
    http://localhost:$query_port/query/service
