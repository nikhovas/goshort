#!/bin/bash


docker pull nikhovas/goshort:redis
docker run -d -p 80:80 -v $PWD/redis-data:/data nikhovas/goshort:redis