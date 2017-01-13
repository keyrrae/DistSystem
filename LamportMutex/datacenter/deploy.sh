#!/usr/bin/env bash
go build
# TODO: deploy via scp
rm dci1/*
rm dci2/*
rm dci3/*
cp ./datacenter dci1/
cp ./datacenter dci2/
cp ./datacenter dci3/
sed "s/4\",/4\",/g" ./servers.conf > ./dci1/servers.conf
sed "s/4\",/5\",/g" ./servers.conf > ./dci2/servers.conf
sed "s/4\",/6\",/g" ./servers.conf > ./dci3/servers.conf