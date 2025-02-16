#!/bin/bash

DEFAULT_BOTS_PATH=`pwd`"/bots"

EXAMPLES_ID=`uuidgen`
EXAMPLES_NAME='Examples'

echo '\c rglua;
INSERT INTO users(id, name) VALUES(:v1, :v2);' | sudo -u postgres psql  \
    -v v1="'$EXAMPLES_ID'" \
    -v v2="'$EXAMPLES_NAME'"

for path in $DEFAULT_BOTS_PATH/*; do
    filename=${path##*/}
    bot=${filename%.lua}
    id=`uuidgen`
    cat ./deployments/load_bots.sql | sudo -u postgres psql \
        -v v1="'$id'" \
        -v v2="'$bot'" \
        -v v3="'$path'" \
        -v v4="'$EXAMPLES_ID'"
done
