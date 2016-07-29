#!/usr/bin/env bash

# Kills all Docker containers

docker kill $(docker ps -a -q)
docker rm $(docker ps -a -q)
