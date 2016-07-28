#!/usr/bin/env bash

# Kills all Docker containers, removes all containers, and deletes saved container images

docker kill $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker rmi $(docker images -q)
