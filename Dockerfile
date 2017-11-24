FROM alpine:3.6
RUN apk update && apk add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY favicon.ico /app
COPY migrations /app/migrations
COPY views /app/views
COPY dist /app/dist
COPY fider /app

CMD [ "./fider" ]