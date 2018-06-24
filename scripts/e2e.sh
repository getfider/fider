#!/bin/bash

FIDER_CONTAINER="fider_e2e"
PG_CONTAINER="fider_pge2e"
PORT=3000

start_fider () {
  echo "Starting Fider (HOST_MODE: $1)..."
  docker rm -f $FIDER_CONTAINER $PG_CONTAINER || true
  docker run -d -e POSTGRES_USER=fider_e2e -e POSTGRES_PASSWORD=fider_e2e_pw --name $PG_CONTAINER postgres:9.6.8
  docker run --link $PG_CONTAINER:pg waisbrot/wait

  if [[ $CI == '' ]]
  then
    docker run \
      -d \
      -p 3000:3000 \
      -e HOST_MODE=$1 \
      -e DATABASE_URL=postgres://fider_e2e:fider_e2e_pw@$PG_CONTAINER:5432/fider_e2e?sslmode=disable \
      --env-file .env \
      --link $PG_CONTAINER \
      --name $FIDER_CONTAINER getfider/fider:e2e
  else
    docker run \
      -d \
      -p 3000:3000 \
      -e HOST_MODE=$1 \
      -e JWT_SECRET=SOME_RANDOM_TOKEN_JUST_FOR_TESTING \
      -e AUTH_ENDPOINT=http://login.dev.fider.io:3000 \
      -e EMAIL_NOREPLY=noreply@fider.io \
      -e EMAIL_MAILGUN_API=$EMAIL_MAILGUN_API \
      -e EMAIL_MAILGUN_DOMAIN=$EMAIL_MAILGUN_DOMAIN \
      -e OAUTH_FACEBOOK_APPID=$OAUTH_FACEBOOK_APPID \
      -e OAUTH_FACEBOOK_SECRET=$OAUTH_FACEBOOK_SECRET \
      -e GO_ENV=development \
      -e DATABASE_URL=postgres://fider_e2e:fider_e2e_pw@$PG_CONTAINER:5432/fider_e2e?sslmode=disable \
      --link $PG_CONTAINER \
      --name $FIDER_CONTAINER getfider/fider:e2e
  fi
}

run_e2e () {
  start_fider $1
  echo "Running e2e tests ..."
  npx jest ./e2e/$1.spec.ts
  if [[ $? == 1 ]] 
  then
    exit 1
  fi
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

echo "Killing Chromium..."
ps -A | grep '[c]hromium' | awk '{print $1}' | xargs kill || true