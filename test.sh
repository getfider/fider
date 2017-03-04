AUTH_ENDPOINT=http://login.test.canhearyou.com:3000 \
GO_ENV=test \
go test $(go list ./... | grep -v /vendor/) -cover