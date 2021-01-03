#!/bin/sh

scope=$1
coll=$2
doc=$3
key=$4

source config.sh

query="statement=insert into \`$COUCHBASE_BUCKET\`.\`$scope\`.\`$coll\` (key,value) values('$key',$doc);"

echo $query

curl -u $COUCHBASE_USER:$COUCHBASE_PASSWORD \
    -d "$query" http://localhost:6093/query/service


