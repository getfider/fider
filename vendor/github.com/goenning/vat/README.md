# Golang VAT number validation

Uses the official [VIES VAT number validation SOAP web service](http://ec.europa.eu/taxation_customs/vies/vatRequest.html?locale=en)
to validate european VAT numbers.

Unfortunately their service is super unreliable.

## Install

```
go get -u github.com/goenning/vat
```


## Usage with Go

```go
import "github.com/goenning/vat"

response, err := vat.Query("IE6388047V")
if err != nil {
  // do sth with err
}
if response.IsValid {
    fmt.Println(response.Name)
}
```

## Retry

Because VIES service is very unreliable, this package will retry every failed request for a maximum of 4 times using an exponential backoff strategy.

## Greece and United Kingdom

This package always assume that country code is on ISO 3166 list.

Greece's ISO 3166 country code is `GR`, but EU uses `EL`. This package does the conversion so that you don't need to care about EU's country code. Same applies to `UK`, both ISO and EU uses `GB`, but this package will always convert `UK` to `GB` for safety reasons.

This conversion only happens for the CountryCode field and not VAT Number. So a EL094160738 is a **valid** VAT number for the `GR` country.

## Fork

This package is a merge of [mattes/vat](https://github.com/mattes/vat) and [dannyvankooten/vat](https://github.com/dannyvankooten/vat), with some refactoring and other features. Thanks!