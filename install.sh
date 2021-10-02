#!/bin/bash

GOBIN=/usr/local/bin/ go install
cp other/goshort.service /etc/systemd/system