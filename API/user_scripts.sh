#!/bin/bash


set -e

cmd="$1"

case "$cmd" in

    create)
        curl "http://localhost:9090/users" \
        -X "POST" \
        -H "Content-type: application/json" \
        -d '{"email":"mizizov@gmail.com","password":"671716Mi","username":"kalzemic", "roles":["Student","User"]}' 
        printf "\n"
        
        ;;
    login)
        curl "http://localhost:9090/users/mizizov@gmail.com?password=671716Mi" \
        -X "GET"
        printf "\n"
        ;;
    get-all)
        curl "http://localhost:9090/users?page=1&size=10" \
        -X "GET"
        printf "\n"
        ;;
    get-domain)
        curl "http://localhost:9090/users?criteria=byEmailDomain&value=gmail.com&page=1&size=10" \
        -X "GET"
        printf "\n"
        ;;
    get-roles)
        curl "http://localhost:9090/users?criteria=byRole&value=User&page=1&size=10" \
        -X "GET"
        printf "\n"
        ;;
    get-reg)
        curl "http://localhost:9090/users?criteria=byRegistrationToday&page=1&size=10" \
        -X "GET"
        printf "\n"
        ;;
    delete)
        curl "http://localhost:9090/users" \
        -X "DELETE"
        printf "\n"
        ;;
    *)
        ;;





esac 