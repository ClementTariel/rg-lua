#!/bin/bash

# TODO: load user and password from conf or env
v1="'referee_temporary_password'"
v2="'matchmaker_temporary_password'"
v3="'rglua_temporary_password'"
cat init_db.sql \
| sed "s/:v1/$v1/g"\
| sed "s/:v2/$v2/g"\
| sed "s/:v3/$v3/g" > init.sql

EXAMPLES_ID=`uuidgen`
EXAMPLES_NAME='Examples'

echo "INSERT INTO users(id, name) VALUES('$EXAMPLES_ID', '$EXAMPLES_NAME');" >> init.sql

LOCAL_BOTS_PATH="../bots"
DOCKER_BOTS_FOLDER="/bots"

for path in $LOCAL_BOTS_PATH/public/* $LOCAL_BOTS_PATH/private/*; do
    filename=${path##*/}
    bot=${filename%.lua}
    path_in_docker="$DOCKER_BOTS_FOLDER/$filename"
    id=`uuidgen`
    echo "INSERT INTO bots(id, name, script, userId, userName) " \
    "VALUES('$id', '$bot', load_bot(CAST('$path_in_docker' AS TEXT)), '$EXAMPLES_ID', '$EXAMPLES_NAME');" >> init.sql
done
