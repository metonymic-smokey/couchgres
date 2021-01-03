#!/bin/sh

source config.sh

scope=$1
coll=$2

curl -u $COUCHBASE_USER:$COUCHBASE_PASSWORD -X POST     \
        http://localhost:6091/pools/default/buckets/$COUCHBASE_BUCKET/collections/$scope \
        -d name=$coll 

