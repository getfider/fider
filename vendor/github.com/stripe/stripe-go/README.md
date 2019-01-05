# Go Stripe

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/stripe/stripe-go)
[![Build Status](https://travis-ci.org/stripe/stripe-go.svg?branch=master)](https://travis-ci.org/stripe/stripe-go)
[![Coverage Status](https://coveralls.io/repos/github/stripe/stripe-go/badge.svg?branch=master)](https://coveralls.io/github/stripe/stripe-go?branch=master)

## Summary

The official [Stripe][stripe] Go client library.

## Versioning

Each revision of the binding is tagged and the version is updated accordingly.

Given Go's lack of built-in versioning, it is highly recommended you use a
[package management tool][package-management] in order to ensure a newer
version of the binding does not affect backwards compatibility.

To see the list of past versions, run `git tag`. To manually get an older
version of the client, clone this repo, checkout the specific tag and build the
library:

```sh
git clone https://github.com/stripe/stripe-go.git
cd stripe-go
git checkout api_version_tag
make build
```

For more details on changes between versions, see the [binding
changelog](CHANGELOG.md) and [API changelog][api-changelog].

## Installation

```sh
go get github.com/stripe/stripe-go
```

## Documentation

For a comprehensive list of examples, check out the [API
documentation][api-docs].

For details on all the functionality in this library, see the [GoDoc][godoc]
documentation.

Below are a few simple examples:

### Customers

```go
params := &stripe.CustomerParams{
	Balance:     stripe.Int64(-123),
	Description: stripe.String("Stripe Developer"),
	Email:       stripe.String("gostripe@stripe.com"),
}
params.SetSource("tok_1234")

customer, err := customer.New(params)
```

### Charges

```go
params := &stripe.ChargeListParams{Customer: stripe.String(customer.ID)}
params.Filters.AddFilter("include[]", "", "total_count")

// set this so you can easily retry your request in case of a timeout
params.Params.IdempotencyKey = stripe.NewIdempotencyKey()

i := charge.List(params)
for i.Next() {
	charge := i.Charge()
}

if err := i.Err(); err != nil {
	// handle
}
```

### Events

```go
i := event.List(nil)
for i.Next() {
	e := i.Event()

	// access event data via e.GetObjectValue("resource_name_based_on_type", "resource_property_name")
	// alternatively you can access values via e.Data.Object["resource_name_based_on_type"].(map[string]interface{})["resource_property_name"]

	// access previous attributes via e.GetPreviousValue("resource_name_based_on_type", "resource_property_name")
	// alternatively you can access values via e.Data.PrevPreviousAttributes["resource_name_based_on_type"].(map[string]interface{})["resource_property_name"]
}
```

Alternatively, you can use the `event.Data.Raw` property to unmarshal to the
appropriate struct.

### Authentication with Connect

There are two ways of authenticating requests when performing actions on behalf
of a connected account, one that uses the `Stripe-Account` header containing an
account's ID, and one that uses the account's keys. Usually the former is the
recommended approach. [See the documentation for more information][connect].

To use the `Stripe-Account` approach, use `SetStripeAccount()` on a `ListParams`
or `Params` class. For example:

```go
// For a list request
listParams := &stripe.ChargeListParams{}
listParams.SetStripeAccount("acct_123")
```

```go
// For any other kind of request
params := &stripe.CustomerParams{}
params.SetStripeAccount("acct_123")
```

To use a key, pass it to `API`'s `Init` function:

```go

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

stripe := &client.API{}
stripe.Init("access_token", nil)
```

### Google AppEngine

If you're running the client in a Google AppEngine environment, you'll need to
create a per-request Stripe client since the `http.DefaultClient` is not
available. Here's a sample handler:

```go
import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

func handler(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        httpClient := urlfetch.Client(c)

        sc := stripeClient.New("sk_live_key", stripe.NewBackends(httpClient))

        chargeParams := &stripe.ChargeParams{
            Amount:      stripe.Int64(2000),
            Currency:    stripe.String(string(stripe.CurrencyUSD)),
            Description: stripe.String("Charge from Google App Engine"),
        }
        chargeParams.SetSource("tok_amex") // obtained with Stripe.js
        charge, err := sc.Charges.New(chargeParams)
        if err != nil {
            fmt.Fprintf(w, "Could not process payment: %v", err)
        }
        fmt.Fprintf(w, "Completed payment: %v", charge.ID)
}
```

## Usage

While some resources may contain more/less APIs, the following pattern is
applied throughout the library for a given `$resource$`:

### Without a Client

If you're only dealing with a single key, you can simply import the packages
required for the resources you're interacting with without the need to create a
client.

```go
import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/$resource$"
)

// Setup
stripe.Key = "sk_key"

stripe.SetBackend("api", backend) // optional, useful for mocking

// Create
$resource$, err := $resource$.New(stripe.$Resource$Params)

// Get
$resource$, err := $resource$.Get(id, stripe.$Resource$Params)

// Update
$resource$, err := $resource$.Update(stripe.$Resource$Params)

// Delete
resourceDeleted, err := $resource$.Del(id, stripe.$Resource$Params)

// List
i := $resource$.List(stripe.$Resource$ListParams)
for i.Next() {
	$resource$ := i.$Resource$()
}

if err := i.Err(); err != nil {
	// handle
}
```

### With a Client

If you're dealing with multiple keys, it is recommended you use `client.API`.
This allows you to create as many clients as needed, each with their own
individual key.

```go
import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

// Setup
sc := &client.API{}
sc.Init("sk_key", nil) // the second parameter overrides the backends used if needed for mocking

// Create
$resource$, err := sc.$Resource$s.New(stripe.$Resource$Params)

// Get
$resource$, err := sc.$Resource$s.Get(id, stripe.$Resource$Params)

// Update
$resource$, err := sc.$Resource$s.Update(stripe.$Resource$Params)

// Delete
resourceDeleted, err := sc.$Resource$s.Del(id, stripe.$Resource$Params)

// List
i := sc.$Resource$s.List(stripe.$Resource$ListParams)
for i.Next() {
	resource := i.$Resource$()
}

if err := i.Err(); err != nil {
	// handle
}
```

### Writing a Plugin

If you're writing a plugin that uses the library, we'd appreciate it if you
identified using `stripe.SetAppInfo`:

```go
stripe.SetAppInfo(&stripe.AppInfo{
    Name:    "MyAwesomePlugin",
    URL:     "https://myawesomeplugin.info",
    Version: "1.2.34",
})
```

This information is passed along when the library makes calls to the Stripe
API. Note that while `Name` is always required, `URL` and `Version` are
optional.

## Development

Pull requests from the community are welcome. If you submit one, please keep
the following guidelines in mind:

1. Code must be `go fmt` compliant.
2. All types, structs and funcs should be documented.
3. Ensure that `make test` succeeds.

## Test

The test suite needs testify's `require` package to run:

    github.com/stretchr/testify/require

Before running the tests, make sure to grab all of the package's dependencies:

    go get -t -v

It also depends on [stripe-mock][stripe-mock], so make sure to fetch and run it from a
background terminal ([stripe-mock's README][stripe-mock-usage] also contains
instructions for installing via Homebrew and other methods):

    go get -u github.com/stripe/stripe-mock
    stripe-mock

Run all tests:

    go test ./...

Run tests for one package:

    go test ./invoice

Run a single test:

    go test ./invoice -run TestInvoiceGet

For any requests, bug or comments, please [open an issue][issues] or [submit a
pull request][pulls].

[api-docs]: https://stripe.com/docs/api/go
[api-changelog]: https://stripe.com/docs/upgrades
[connect]: https://stripe.com/docs/connect/authentication
[godoc]: http://godoc.org/github.com/stripe/stripe-go
[issues]: https://github.com/stripe/stripe-go/issues/new
[package-management]: https://code.google.com/p/go-wiki/wiki/PackageManagementTools
[pulls]: https://github.com/stripe/stripe-go/pulls
[stripe]: https://stripe.com
[stripe-mock]: https://github.com/stripe/stripe-mock
[stripe-mock-usage]: https://github.com/stripe/stripe-mock#usage

<!--
# vim: set tw=79:
-->
