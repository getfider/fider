#####################
### Server Build Step
#####################
FROM --platform=${TARGETPLATFORM} golang:1.17-buster AS server-builder 

ARG BUILDNUMBER=local
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN mkdir /server
WORKDIR /server

COPY . .
RUN BUILDNUMBER=${BUILDNUMBER} GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build-server

#################
### UI Build Step
#################
FROM --platform=${TARGETPLATFORM} node:16-buster AS ui-builder 

WORKDIR /ui

COPY . .
RUN npm ci
RUN make build-ssr
RUN make build-ui

################
### Runtime Step
################
FROM --platform=${TARGETPLATFORM} debian:buster-slim

RUN apt-get update
RUN apt-get install -y ca-certificates

WORKDIR /app

COPY --from=server-builder /server/migrations /app/migrations
COPY --from=server-builder /server/views /app/views
COPY --from=server-builder /server/locale /app/locale
COPY --from=server-builder /server/LICENSE /app
COPY --from=server-builder /server/fider /app

COPY --from=ui-builder /ui/favicon.png /app
COPY --from=ui-builder /ui/dist /app/dist
COPY --from=ui-builder /ui/robots.txt /app
COPY --from=ui-builder /ui/ssr.js /app

EXPOSE 3000

HEALTHCHECK --timeout=5s CMD ./fider ping

CMD ./fider migrate && ./fider