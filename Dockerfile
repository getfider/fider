FROM alpine:3.4
RUN apk add --no-cache ca-certificates

WORKDIR /

COPY wechy /
COPY migrations /migrations
COPY views /views
COPY dist /dist
COPY favicon.ico /

ENV PORT 8080
EXPOSE 8080

CMD ["/wechy"]