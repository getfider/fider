FROM alpine:3.6
RUN apk update && apk add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY favicon.ico /app
COPY migrations /app/migrations
COPY views /app/views
COPY dist /app/dist
COPY LICENSE /app
COPY fider /app

EXPOSE 3000

HEALTHCHECK --timeout=5s CMD ./fider ping

CMD ./fider migrate

CMD [ "./fider" ]