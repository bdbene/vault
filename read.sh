#!/bin/bash

for ((i=0; i<=1000; i++)); do
    curl -s -H "Content-Type: application/json" --request POST --data '{"identifier": "curl'"${1}${i}"'", "password": "curlPass"}' --url http://localhost:8080/secret/query
    echo ""
done
