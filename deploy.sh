#!/usr/bin/env bash

git pull
./build.sh
docker stack deploy --with-registry-auth --resolve-image=always -c ./docker-stack.yml leto