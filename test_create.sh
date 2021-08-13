#!/bin/sh

for i in {1..5}; do
  STATUSCODE=$(curl --silent --header "Content-Type: application/json"  --output /dev/stderr --write-out "%{http_code}" -d '{"title": "Test todo '"$i"'"}' ${CREATE_URL})
  if test $STATUSCODE -ne 200; then
    exit 1
  fi
done
exit 0
