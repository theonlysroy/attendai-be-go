#!/usr/bin/env bash

url="http://localhost:4040/api"

# health check
echo -e "Checking health...\n"
curl -sS "${url}/healthz" | jq

# get students
echo -e "Get all students...\n"
curl -sS "${url}/students" | jq
