#!/bin/bash

set -e

cmd="$1"

case $cmd in

    poll)
    curl "http://localhost:9094/posts/poll" \
    -X "POST" \
    -H "Content-type: application/json" \
    -d '{"maxPosts":10,"maxDuration":100000000}'
    printf "\n"
    ;;
    *)
    ;;

esac