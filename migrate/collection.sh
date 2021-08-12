#!/bin/sh

source config.sh

scope=$1
coll=$2

curl -u $COUCHBASE_USER:$COUCHBASE_PASSWORD -X POST     \
        http://localhost:$COUCHBASE_PORT/pools/default/buckets/$COUCHBASE_BUCKET/scopes/$scope/collections \
        -d name=$coll 
