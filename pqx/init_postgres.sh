#!/usr/bin/env bash

db=$1
dockerName=postgres5432
port=5432

if [ -z "$db" ]
then
    db="db"
fi

echo "Remove old container..."
docker rm -f -v $dockerName

echo "Create new postgres docker container..."
docker run -d --name $dockerName -p $port:5432 -e POSTGRES_USER=postgres postgres:latest

echo "Checking Postgres availability..."
while : ; do
    echo "Waiting for Postgres..."
    sleep 3
    upstartCheck=$(docker exec $dockerName /bin/sh -c "ps aux | grep 'postgres' | grep 'docker-entrypoint.sh' | grep -v 'grep'")
    [[ $upstartCheck == *"docker-entrypoint.sh"* ]] || break
done

echo "Create database..."
docker exec $dockerName /bin/sh -c "su postgres --command 'createdb -O postgres $db'"
