#!/bin/bash


redis-server &
P1=$!
goshort &
P2=$!
wait $P1 $P2