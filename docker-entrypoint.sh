#!/bin/bash

set -e

if [[ "$1" == "app" ]]; then
    exec linuxmender
elif [[ "$1" == "test" ]]; then
    exec go test ./...
else
    exec "$@"
fi
