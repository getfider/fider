#####################
### Server Build Step
#####################
FROM --platform=${TARGETPLATFORM:-linux/amd64} golang:1.22-bullseye AS server-builder 


RUN mkdir /server
WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG COMMITHASH
RUN COMMITHASH=${COMMITHASH} GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build-server

#################
### UI Build Step
#################
FROM --platform=${TARGETPLATFORM:-linux/amd64} node:18-bullseye AS ui-builder 

WORKDIR /ui

COPY package.json package-lock.json ./
RUN npm ci --maxsockets 1

COPY . .
RUN make build-ssr
RUN make build-ui

################
### Runtime Step
################
FROM --platform=${TARGETPLATFORM:-linux/amd64} debian:bullseye-slim

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