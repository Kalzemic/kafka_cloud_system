#!/bin/bash 

set -e

docker compose up -d

echo "Waiting for Kafka..."

sleep 7
echo "Kafka ready"


for dir in */; do 
    (
        cd "$dir" || exit
        if [[ "$dir" != "client/" ]]; then
            go mod download
            go run .
        else
            [ -d node_modules ] || npm install
            npm run dev
        fi
    ) &
done

wait