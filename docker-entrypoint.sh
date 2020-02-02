#!/bin/bash

if [[ "$1" == 'app' ]]; then
    ./linuxmender
elif [[ "$1" == 'test' ]]; then
    go test ./...
fi
