#!/usr/bin/env bash

docker run --restart=always -d -p 80:80 byuoitav/raspi-tp:latest
docker run --restart=always -d -p 8000:8000 byuoitav/rpi-av-api:latest
docker run --restart=always -d -p 8005:8005 byuoitav/rpi-pjlink-microservice:latest
docker run --restart=always -d -p 8007:8007 byuoitav/rpi-sony-control-microservice:latest
docker run --restart=always -d -p 8001:8001 byuoitav/rpi-telnet-microservice:latest
