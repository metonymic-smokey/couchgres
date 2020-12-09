#!/bin/sh

source config.sh

curl -X GET -u $COUCHBASE_USERNAME:$COUCHBASE_PASSWORD http://127.0.0.1:$COUCHBASE_PORT/pools/default/buckets/$COUCHBASE_BUCKET/collections | jq '.' > scopes.json 

# probably by running queries concurrently(serially to start with) to find the number of items in each scope/collection.
