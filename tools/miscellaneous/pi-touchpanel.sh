#!/usr/bin/env bash

# Sets up local AV API Docker containers on a Raspberry Pi touchpanel (run after drydock.sh for best effect)

docker pull byuoitav/rpi-av-api:latest
docker kill av-api
docker rm av-api
docker run -e PI_TOUCHPANEL="true" -e EMS_API_USERNAME=$EMS_API_USERNAME -e EMS_API_PASSWORD=$EMS_API_PASSWORD -d --restart=always --name av-api -p 8000:8000 byuoitav/rpi-av-api:latest

docker pull byuoitav/rpi-telnet-microservice:latest
docker kill telnet-microservice
docker rm telnet-microservice
docker run -e PI_TOUCHPANEL="true" -d --restart=always --name telnet-microservice -p 8001:8001 byuoitav/rpi-telnet-microservice:latest

docker pull byuoitav/rpi-pjlink-microservice:latest
docker kill pjlink-service
docker rm pjlink-service
docker run -e PI_TOUCHPANEL="true" -d --restart=always --name pjlink-service -p 8005:8005 byuoitav/rpi-pjlink-microservice:latest

docker pull byuoitav/rpi-configuration-database-microservice:latest
docker kill configuration-database-microservice
docker rm configuration-database-microservice
docker run -e PI_TOUCHPANEL="true" -e CONFIGURATION_DATABASE_USERNAME=$CONFIGURATION_DATABASE_USERNAME -e CONFIGURATION_DATABASE_PASSWORD=$CONFIGURATION_DATABASE_PASSWORD -e CONFIGURATION_DATABASE_HOST=$CONFIGURATION_DATABASE_HOST -e CONFIGURATION_DATABASE_PORT=$CONFIGURATION_DATABASE_PORT -e CONFIGURATION_DATABASE_NAME=$CONFIGURATION_DATABASE_NAME -d --restart=always --name configuration-database-microservice -p 8006:8006 byuoitav/rpi-sony-control-microservice:latest

docker pull byuoitav/rpi-sony-control-microservice:latest
docker kill sony-control-microservice
docker rm sony-control-microservice
docker run -e PI_TOUCHPANEL="true" -e SONY_TV_PSK=$SONY_TV_PSK -d --restart=always --name sony-control-microservice -p 8007:8007 byuoitav/rpi-sony-control-microservice:latest

docker pull byuoitav/raspi-tp:latest
docker kill raspi-tp
docker rm raspi-tp
docker run -e PI_TOUCHPANEL="true" --restart=always -d --name raspi-tp -p 8888:8888 byuoitav/raspi-tp:latest
