#!/usr/bin/env bash

# Sets up local AV API Docker containers on a Raspberry Pi touchpanel (run after drydock.sh for best effect)

docker pull byuoitav/raspi-tp:latest
docker pull byuoitav/rpi-av-api:latest
docker pull byuoitav/rpi-telnet-microservice:latest
docker pull byuoitav/rpi-pjlink-microservice:latest
docker pull byuoitav/rpi-sony-control-microservice:latest

docker run -e PI_TOUCHPANEL="true" --restart=always -d -p 80:80 byuoitav/raspi-tp:latest
docker run -e PI_TOUCHPANEL="true" -e EMS_API_USERNAME=$EMS_API_USERNAME -e EMS_API_PASSWORD=$EMS_API_PASSWORD -d --restart=always --name av-api -p 8000:8000 byuoitav/rpi-av-api:latest
docker run -e PI_TOUCHPANEL="true" -d --restart=always --name telnet-microservice -p 8001:8001 byuoitav/rpi-telnet-microservice:latest
docker run -e PI_TOUCHPANEL="true" -d --restart=always --name pjlink-service -p 8005:8005 byuoitav/rpi-pjlink-microservice:latest
docker run -e PI_TOUCHPANEL="true" -e SONY_TV_PSK=$SONY_TV_PSK -d --restart=always --name sony-control-microservice -p 8007:8007 byuoitav/rpi-sony-control-microservice:latest
