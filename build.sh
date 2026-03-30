#!/bin/bash

echo "building..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/app .

if [[ $1 == "d" ]]; then
    echo "deploying..."
    cd ../../..
    docker compose restart learn-jap-v6
fi

echo "done!"