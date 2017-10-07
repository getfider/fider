#!/bin/bash

echo 'mode: atomic' > cover.out && go list ./... | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> cover.out' && rm coverage.tmp
