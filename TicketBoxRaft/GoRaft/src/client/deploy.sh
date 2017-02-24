#!/usr/bin/env bash

go build

cp -f ./client cli1/
cp -f ./client cli2/
cp -f ./client cli3/
sed "s/4\",/4\",/g" ./client.conf > ./cli1/client.conf
sed "s/4\",/5\",/g" ./client.conf > ./cli2/client.conf
sed "s/4\",/6\",/g" ./client.conf > ./cli3/client.conf

