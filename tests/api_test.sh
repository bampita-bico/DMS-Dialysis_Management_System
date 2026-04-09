#! /bin/bash

echo "Testing API..."

responses=$(curl -s https://localhost:8080/health)

if [[ $response == *"ok"* ]]; then
  echo "Health endpoint working"
else
  echo "Health endpoint failed"
  exit 1
fi
