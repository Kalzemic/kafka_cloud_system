#!/bin/bash


set -e

cmd="$1"

case "$cmd" in

    create)
        curl "http://localhost:9091/users" \
        -X "POST" \
        -H "Content-type: application/json" \
        -d '{"email":"mizizov@gmail.com","password":"671716Mi","username":"kalzemic", "roles":["Student","User"]}'
        ;;
    login)
        curl "http://localhost:9091/users/mizizov@gmail.com?password=671716Mi" \
        -X "GET"
        ;;
    get-all)
        curl "http://localhost:9091/users?page=1&size=10" \
        -X "GET"
        ;;
    get-domain)
        curl "http://localhost:9091/users?criteria=byEmailDomain&value=gmail.com&page=1&size=10" \
        -X "GET"
        ;;
    get-roles)
        curl "http://localhost:9091/users?criteria=byRole&value=User&page=1&size=10" \
        -X "GET"
        ;;
    get-reg)
        curl "http://localhost:9091/users?criteria=byRegistrationToday&page=1&size=10" \
        -X "GET"
        ;;
    delete)
        curl "http://localhost:9091/users" \
        -X "DELETE"
        ;;
    *)
        ;;





esac 