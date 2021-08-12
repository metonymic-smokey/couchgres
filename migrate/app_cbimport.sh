#!/bin/sh

source config.sh

scope=$1
coll=$2
file=$3
key=$4

PORT=6091

./cbimport csv -c http://localhost:$PORT -u "$COUCHBASE_USER" -p "$COUCHBASE_PASSWORD" -b "$COUCHBASE_BUCKET" --generate-key %$key%::#MONO_INCR# -d "file://$file" --scope-collection-exp $scope.$coll

