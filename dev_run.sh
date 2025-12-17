#!/bin/bash 

set -e

docker compose up -d


for dir in */; do 
    (
        cd "$dir" || exit
        go run .
    ) &
done

wait