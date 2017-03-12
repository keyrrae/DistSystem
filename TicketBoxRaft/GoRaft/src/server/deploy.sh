#!/usr/bin/env bash
go build
# TODO: deploy via scp
rm dci1/*
rm dci2/*
rm dci3/*
rm dci4/*
rm dci5/*
cp ./server dci1/
cp ./server dci2/
cp ./server dci3/
cp ./server dci4/
cp ./server dci5/
sed "s/\"self\": \"127.0.0.1:1234\",/\"self\": \"127.0.0.1:1234\",/g" ./server_conf.json | sed "s/ 1,/ 1,/g" > ./dci1/server_conf.json
sed "s/\"self\": \"127.0.0.1:1234\",/\"self\": \"127.0.0.1:1235\",/g" ./server_conf.json | sed "s/ 1,/ 2,/g" > ./dci2/server_conf.json
sed "s/\"self\": \"127.0.0.1:1234\",/\"self\": \"127.0.0.1:1236\",/g" ./server_conf.json | sed "s/ 1,/ 3,/g" > ./dci3/server_conf.json

sed "s/\"self\": \"127.0.0.1:1237\"/\"self\": \"127.0.0.1:1237\"/g" ./new_conf.json | sed "s/\"processid\": 4/\"processid\": 4/g" > ./dci4/server_conf.json
sed "s/\"self\": \"127.0.0.1:1237\"/\"self\": \"127.0.0.1:1238\"/g" ./new_conf.json | sed "s/\"processid\": 4/\"processid\": 5/g" > ./dci5/server_conf.json

