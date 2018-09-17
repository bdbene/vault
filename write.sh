#!/bin/bash

for ((i=0; i<=5000; i++)); do
    curl -s -H "Content-Type: application/json" --request POST --data '{"identifier": "curl'"${1}${i}"'", "password": "curlPass", "secret": "SentWithCurl"}' --url http://localhost:8080/secret/write &
    echo ""
done
