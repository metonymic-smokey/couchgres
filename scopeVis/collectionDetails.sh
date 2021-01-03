#!/bin/sh

source config.sh

docker exec -it $COUCHBASE_CONTAINER bash -c "/opt/couchbase/bin/cbstats localhost:11210 -u $COUCHBASE_USERNAME -p $COUCHBASE_PASSWORD -b $COUCHBASE_BUCKET collections" > details

#stores scope and collection details in a file called 'details'
