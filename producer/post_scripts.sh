#!/bin/bash

set -e


cmd="$1"

case $cmd in 

    produce)
    curl "http://localhost:9093/posts/produce" \
    -X "POST" \
    -H "Content-type: application/json" \
    -d '{"email":"mizizov@gmail.com","content":"welcome from cloud computing class","timestamp":"2025-12-17T10:42:31.527Z"}'
    printf "\n"
    ;;

    *)
    ;;

esac