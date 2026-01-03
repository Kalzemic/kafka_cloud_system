#!/bin/bash 

set -e

docker compose up -d


for dir in */; do 
    (
        cd "$dir" || exit
        if [[ "$dir" != "client/" ]]; then
            go run .
        else
            [ -d node_modules ] || npm install
            npm run dev
        fi
    ) &
done

wait