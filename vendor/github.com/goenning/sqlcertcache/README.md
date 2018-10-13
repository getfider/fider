[![GoDoc](https://godoc.org/github.com/goenning/sqlcertcache?status.svg)](https://godoc.org/github.com/goenning/sqlcertcache) [![Go Report Card](https://goreportcard.com/badge/github.com/goenning/sqlcertcache)](https://goreportcard.com/report/github.com/goenning/sqlcertcache)

# sqlcertcache

SQL cache for [acme/autocert](https://godoc.org/golang.org/x/crypto/acme/autocert) written in Go.

## Example

```go
conn, err := sql.Open("postgres", "postgres://YOUR_CONNECTION_STRING")
if err != nil {
  // Handle error
}

cache, err := sqlcertcache.New(conn, "autocert_cache")
if err != nil {
  // Handle error
}

m := autocert.Manager{
  Prompt:     autocert.AcceptTOS,
  HostPolicy: autocert.HostWhitelist("example.org"),
  Cache:      cache,
}

s := &http.Server{
  Addr:      ":https",
  TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
}

s.ListenAndServeTLS("", "")
```

## Performance

This is just a reminder that autocert has an internal in-memory cache that is used before quering this long-term cache.
So you don't need to worry about your SQL database being hit many times just to get a certificate. It should only do once per process+key.

## Thanks

Inspired by https://github.com/danilobuerger/autocert-s3-cache

## License

MIT