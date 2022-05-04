#!/bin/bash
PROJECTDIR=$(git rev-parse --show-toplevel)
cd $PROJECTDIR/src
# Unless already defined
if [ -z $SNOOPY_VESRION ]; then
    export SNOOPY_VERSION=v0.1.1
fi
docker build --tag dfroberg/snoopy:$SNOOPY_VERSION .
docker image tag dfroberg/snoopy:$SNOOPY_VERSION dfroberg/snoopy:$SNOOPY_VERSION
docker image tag dfroberg/snoopy:$SNOOPY_VERSION dfroberg/snoopy:latest
docker image push --all-tags dfroberg/snoopy