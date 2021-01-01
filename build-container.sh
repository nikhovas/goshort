#!/usr/bin/env bash

docker build . -t nikhovas/goshort:alpine
docker push nikhovas/goshort:alpine