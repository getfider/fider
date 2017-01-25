FROM golang:1.7.4
RUN go get -u github.com/kardianos/govendor 

ADD . /go/src/github.com/WeCanHearYou/wchy-api
WORKDIR /go/src/github.com/WeCanHearYou/wchy-api

RUN govendor sync
RUN make build

ENV PORT 8080

ENTRYPOINT ./wchy-api

EXPOSE 8080