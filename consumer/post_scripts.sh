#!/bin/bash

set -e

cmd="$1"

case $cmd in


    listen)
        curl "http://localhost:9094/posts/listen" \
        -X "GET"
        printf "\n"
    ;;
    *)
    ;;

esac