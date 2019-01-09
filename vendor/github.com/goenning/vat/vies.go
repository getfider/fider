package vat

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const serviceURL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"

//DefaultClient is the HTTP Client used to communicate with VIES
var DefaultClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
}

var envelope *template.Template

var rd struct {
	XMLName xml.Name `xml:"Envelope"`
	Soap    struct {
		XMLName xml.Name `xml:"Body"`
		Soap    struct {
			XMLName     xml.Name `xml:"checkVatResponse"`
			CountryCode string   `xml:"countryCode"`
			VATNumber   string   `xml:"vatNumber"`
			Valid       bool     `xml:"valid"`
			Name        string   `xml:"name"`
			Address     string   `xml:"address"`
		}
		SoapFault struct {
			XMLName string `xml:"Fault"`
			Code    string `xml:"faultcode"`
			Message string `xml:"faultstring"`
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	t, err := template.New("envelope").Parse(`
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://schemas.conversesolutions.com/xsd/dmticta/v1">
	<soapenv:Header/>
	<soapenv:Body>
	  <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
		<countryCode>{{.countryCode}}</countryCode>
		<vatNumber>{{.vatNumber}}</vatNumber>
	  </checkVat>
	</soapenv:Body>
	</soapenv:Envelope>
	`)
	if err != nil {
		panic("vat: could not parse VIES envelope")
	}
	envelope = t
}

func retry(attempts int, sleep time.Duration, f func() bool) {
	shouldRetry := f()
	attempts--
	if shouldRetry && attempts > 0 {
		time.Sleep(sleep)
		retry(attempts, 2*sleep, f)
	}
}

func sendRequestVIES(vatNumber string) (*Response, error) {
	var body bytes.Buffer
	if err := envelope.Execute(&body, map[string]string{
		"countryCode": strings.ToUpper(vatNumber[0:2]),
		"vatNumber":   vatNumber[2:],
	}); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", serviceURL, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("Connection", "close")

	res, err := DefaultClient.Do(req)
	if err != nil {
		return nil, ErrServiceUnreachable
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, ErrVIESServiceUnavailable
	}

	xmlResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(xmlResponse, &rd); err != nil {
		return nil, err
	}

	if rd.Soap.SoapFault.Message != "" {
		if rd.Soap.SoapFault.Message == "INVALID_INPUT" {
			return nil, ErrVIESInvalidInput
		} else if rd.Soap.SoapFault.Message == "GLOBAL_MAX_CONCURRENT_REQ" {
			return nil, ErrVIESGlobalMaxConcurrentRequest
		} else if rd.Soap.SoapFault.Message == "MS_MAX_CONCURRENT_REQ" {
			return nil, ErrVIESMSMaxConcurrentRequest
		} else if rd.Soap.SoapFault.Message == "SERVICE_UNAVAILABLE" {
			return nil, ErrVIESServiceUnavailable
		} else if rd.Soap.SoapFault.Message == "MS_UNAVAILABLE" {
			return nil, ErrVIESMSUnavailable
		} else if rd.Soap.SoapFault.Message == "TIMEOUT" {
			return nil, ErrVIESTimeout
		}
	}

	return &Response{
		CountryCode: toISOCountryCode(rd.Soap.Soap.CountryCode),
		VATNumber:   rd.Soap.Soap.VATNumber,
		IsValid:     rd.Soap.Soap.Valid,
		Name:        rd.Soap.Soap.Name,
		Address:     formatAddress(rd.Soap.Soap.Address),
	}, nil
}
