#!/bin/bash

source config.sh

scope=$1
query_port=$(($COUCHBASE_PORT+2))

curl -u $COUCHBASE_USER:$COUCHBASE_PASSWORD \
    -d "statement=create scope \`$COUCHBASE_BUCKET\`.$scope;" \
    http://localhost:$query_port/query/service
