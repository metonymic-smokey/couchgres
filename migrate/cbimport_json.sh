#!/bin/sh

source config.sh

scope=$1
coll=$2
file=$3
key=$4

docker cp $file $CONTAINER_NAME:/$file

cmd="/opt/couchbase/bin/cbimport json -c http://localhost:8091 -u "$COUCHBASE_USER" -p "$COUCHBASE_PASSWORD" -b '$COUCHBASE_BUCKET' -d 'file://$file' -f list --generate-key %$key%::#MONO_INCR# --scope-collection-exp '$scope.$coll'"

docker exec $CONTAINER_NAME /bin/sh -c "$cmd"




