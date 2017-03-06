#!/usr/bin/env bash

go build

cp -f ./client cli1/
cp -f ./client cli2/
cp -f ./client cli3/
sed "s/4\",/4\",/g" ./client_conf.json > ./cli1/client_conf.json
sed "s/4\",/5\",/g" ./client_conf.json > ./cli2/client_conf.json
sed "s/4\",/6\",/g" ./client_conf.json > ./cli3/client_conf.json
