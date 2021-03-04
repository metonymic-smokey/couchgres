#!/bin/sh

source config.sh

scope=$1
coll=$2
file=$3
key=$4

./cbimport json -c http://localhost:$COUCHBASE_PORT -u "$COUCHBASE_USER" -p "$COUCHBASE_PASSWORD" -b $COUCHBASE_BUCKET -d file://$file -f list --generate-key %$key%::#MONO_INCR# --scope-collection-exp $scope.$coll




