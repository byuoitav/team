#!/bin/bash

BEARER_TOKEN=$(get-bearer-token)

curl -H "Authorization: Bearer $BEARER_TOKEN" https://byuoitav-raspi-deployment-microservice.avs.byu.edu/webhook
