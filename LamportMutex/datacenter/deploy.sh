#!/usr/bin/env bash
go build
# TODO: deploy via scp
rm dci1/*
rm dci2/*
rm dci3/*
cp ./datacenter dci1/
cp ./datacenter dci2/
cp ./datacenter dci3/
sed "s/4\",/4\",/g" ./server_conf.json | sed "s/ 1,/ 1,/g" > ./dci1/server_conf.json
sed "s/4\",/5\",/g" ./server_conf.json | sed "s/ 1,/ 2,/g" > ./dci2/server_conf.json
sed "s/4\",/6\",/g" ./server_conf.json | sed "s/ 1,/ 3,/g" > ./dci3/server_conf.json
