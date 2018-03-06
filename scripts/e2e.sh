#!/bin/bash

FIDER_CONTAINER="fider_e2e"
PG_CONTAINER="fider_pge2e"
PORT=3000

run_e2e () {
  # Check is PORT is in use
  if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null ; then
      echo "Another process is already running on port $PORT."
      exit 1;
  fi

  echo "Compiling tests..."
  rm -rf output
  ./node_modules/.bin/tsc -p ./tests/tsconfig.json

  echo "Starting Fider (HOST_MODE: $1)..."
  docker rm -f $PG_CONTAINER || true
  docker run -d -e POSTGRES_USER=fider_e2e -e POSTGRES_PASSWORD=fider_e2e_pw --name $PG_CONTAINER postgres:9.6.2
  docker run --link $PG_CONTAINER:pg waisbrot/wait
  docker run -d -p 3000:3000 --link $PG_CONTAINER -e HOST_MODE=$1 -e DATABASE_URL=postgres://fider_e2e:fider_e2e_pw@$PG_CONTAINER:5432/fider_e2e?sslmode=disable --env-file .env --name $FIDER_CONTAINER getfider/fider:e2e

  {
    {
      echo "Running e2e tests ..."
      ./node_modules/.bin/mocha -t 60000 output/tests/e2e-$1.js
    } || { 
      echo "Tests failed..."; 
    }
  } && {
      echo "Stopping Fider ..."
      docker logs $FIDER_CONTAINER >> ./logs/e2e.log
      docker rm -f $FIDER_CONTAINER
  }
}

if [[ $1 == 'build' ]] || [ -z $1 ]
then
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
  docker build -t getfider/fider:e2e .
fi

if [[ $1 == 'single' ]] || [ -z $1 ]
then
  run_e2e single
fi

if [[ $1 == 'multi' ]] || [ -z $1 ]
then
  run_e2e multi
fi


echo "Stopping Postgres ..."
docker rm -f $PG_CONTAINER || true