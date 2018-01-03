#!/bin/bash

CONTAINER="fider_pge2e"
PORT=3000

if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null ; then
    echo "Another process is already running on port $PORT."
    exit 1;
fi

echo "Clean up ..."
rm -rf output
tsc
mkdir -p logs/
docker rm -f $CONTAINER || true
docker run -d -e POSTGRES_USER=fider_e2e -e POSTGRES_PASSWORD=fider_e2e_pw -p 5577:5432 --name $CONTAINER postgres:9.6.2
docker run --link $CONTAINER -e TARGETS=$CONTAINER:5432 waisbrot/wait

run_e2e () {
  echo "Starting Fider (HOST_MODE: $1)..."
  HOST_MODE=$1 DATABASE_URL=postgres://fider_e2e:fider_e2e_pw@localhost:5577/fider_e2e?sslmode=disable godotenv -f .env ./fider > logs/e2e.log 2>&1 &
  FIDER_PID=$!

  {
    {
      echo "Running e2e tests ..."
      ./node_modules/.bin/mocha -t 60000 output/tsc/tests/e2e-$1.js
    } || { 
      echo "Tests failed..."; 
    }
  } && {
      echo "Killing Fider -> $FIDER_PID ..."
      kill $(ps -ef | grep $FIDER_PID | grep -v grep | awk '{print $2}')
  }
}

run_e2e single
run_e2e multi

echo "Killing Container -> $CONTAINER ..."
docker rm -f $CONTAINER || true