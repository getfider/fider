FROM alpine:3.8
RUN apk update && apk add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY favicon.png /app
COPY migrations /app/migrations
COPY views /app/views
COPY dist /app/dist
COPY LICENSE /app
COPY robots.txt /app
COPY fider /app

EXPOSE 3000

HEALTHCHECK --timeout=5s CMD ./fider ping

CMD [ "./fider" ]