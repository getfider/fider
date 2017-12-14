#!/bin/bash

CONTAINER="fider_pge2e"

# TODO: check if something is running on port 3000 and then abort

echo "Clean up ..."
rm -rf output
tsc
docker rm -f $CONTAINER || true
docker run -d -e POSTGRES_USER=fider_e2e -e POSTGRES_PASSWORD=fider_e2e_pw -p 5577:5432 --name $CONTAINER postgres:9.6.2
docker run --link $CONTAINER -e TARGETS=$CONTAINER:5432 waisbrot/wait

echo "Starting Fider (Single Host)..."
HOST_MODE=single DATABASE_URL=postgres://fider_e2e:fider_e2e_pw@localhost:5577/fider_e2e?sslmode=disable godotenv -f .env ./fider > logs/e2e.log 2>&1 &
FIDER_PID=$!

{
  {
    echo "Running e2e tests ..."
    ./node_modules/.bin/mocha -t 60000 output/tsc/tests/e2e-single.js
  } || { 
    echo "Tests failed..."; 
  }
} && {
    echo "Killing Fider -> $FIDER_PID ..."
    kill $(ps -ef | grep $FIDER_PID | grep -v grep | awk '{print $2}')
}

echo "Starting Fider (Multi Host)..."
HOST_MODE=multi DATABASE_URL=postgres://fider_e2e:fider_e2e_pw@localhost:5577/fider_e2e?sslmode=disable godotenv -f .env ./fider > logs/e2e.log 2>&1 &
FIDER_PID=$!

{
  {
    echo "Running e2e tests ..."
    ./node_modules/.bin/mocha -t 60000 output/tsc/tests/e2e-multi.js
  } || { 
    echo "Tests failed..."; 
  }
} && {
    echo "Killing Fider -> $FIDER_PID ..."
    kill $(ps -ef | grep $FIDER_PID | grep -v grep | awk '{print $2}')
}

echo "Killing Container -> $CONTAINER ..."
docker rm -f $CONTAINER || true