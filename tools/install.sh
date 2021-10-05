#!/bin/bash

GOBIN=/usr/local/bin/ go install
cp tools/supervisor.conf /etc/supervisor/conf.d/goshort.conf