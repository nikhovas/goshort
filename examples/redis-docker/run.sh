#!/bin/bash


docker run -d -p 80:80 -v $PWD/redis-data:/data nikhovas/goshort:redis