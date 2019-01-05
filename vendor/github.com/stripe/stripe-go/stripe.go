// Package stripe provides the binding for Stripe REST APIs.
package stripe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/stripe/stripe-go/form"
)

//
// Public constants
//

const (
	// APIBackend is a constant representing the API service backend.
	APIBackend SupportedBackend = "api"

	// APIURL is the URL of the API service backend.
	APIURL string = "https://api.stripe.com"

	// UnknownPlatform is the string returned as the system name if we couldn't get
	// one from `uname`.
	UnknownPlatform string = "unknown platform"

	// UploadsBackend is a constant representing the uploads service backend.
	UploadsBackend SupportedBackend = "uploads"

	// UploadsURL is the URL of the uploads service backend.
	UploadsURL string = "https://files.stripe.com"
)

//
// Public variables
//

// Key is the Stripe API key used globally in the binding.
var Key string

// LogLevel is the logging level for this library.
// 0: no logging
// 1: errors only
// 2: errors + informational (default)
// 3: errors + informational + debug
var LogLevel = 2

// Logger controls how stripe performs logging at a package level. It is useful
// to customise if you need it prefixed for your application to meet other
// requirements.
//
// This Logger will be inherited by any backends created by default, but will
// be overridden if a backend is created with GetBackendWithConfig.
var Logger Printfer

//
// Public types
//

// AppInfo contains information about the "app" which this integration belongs
// to. This should be reserved for plugins that wish to identify themselves
// with Stripe.
type AppInfo struct {
	Name      string `json:"name"`
	PartnerID string `json:"partner_id"`
	URL       string `json:"url"`
	Version   string `json:"version"`
}

// formatUserAgent formats an AppInfo in a way that's suitable to be appended
// to a User-Agent string. Note that this format is shared between all
// libraries so if it's changed, it should be changed everywhere.
func (a *AppInfo) formatUserAgent() string {
	str := a.Name
	if a.Version != "" {
		str += "/" + a.Version
	}
	if a.URL != "" {
		str += " (" + a.URL + ")"
	}
	return str
}

// Backend is an interface for making calls against a Stripe service.
// This interface exists to enable mocking for during testing if needed.
type Backend interface {
	Call(method, path, key string, params ParamsContainer, v interface{}) error
	CallRaw(method, path, key string, body *form.Values, params *Params, v interface{}) error
	CallMultipart(method, path, key, boundary string, body *bytes.Buffer, params *Params, v interface{}) error
	SetMaxNetworkRetries(maxNetworkRetries int)
}

// BackendConfig is used to configure a new Stripe backend.
type BackendConfig struct {
	// HTTPClient is an HTTP client instance to use when making API requests.
	//
	// If left unset, it'll be set to a default HTTP client for the package.
	HTTPClient *http.Client

	// LogLevel is the logging level of the library and defined by:
	//
	// 0: no logging
	// 1: errors only
	// 2: errors + informational (default)
	// 3: errors + informational + debug
	//
	// Defaults to 0 (no logging), so please make sure to set this if you want
	// to see logging output in your custom configuration.
	LogLevel int

	// Logger is where this backend will write its logs.
	//
	// If left unset, it'll be set to Logger.
	Logger Printfer

	// MaxNetworkRetries sets maximum number of times that the library will
	// retry requests that appear to have failed due to an intermittent
	// problem.
	//
	// Defaults to 0.
	MaxNetworkRetries int

	// URL is the base URL to use for API paths.
	//
	// If left empty, it'll be set to the default for the SupportedBackend.
	URL string
}

// BackendImplementation is the internal implementation for making HTTP calls
// to Stripe.
//
// The public use of this struct is deprecated. It will be unexported in a
// future version.
type BackendImplementation struct {
	Type              SupportedBackend
	URL               string
	HTTPClient        *http.Client
	MaxNetworkRetries int
	LogLevel          int
	Logger            Printfer

	// networkRetriesSleep indicates whether the backend should use the normal
	// sleep between retries.
	//
	// See also SetNetworkRetriesSleep.
	networkRetriesSleep bool
}

// Call is the Backend.Call implementation for invoking Stripe APIs.
func (s *BackendImplementation) Call(method, path, key string, params ParamsContainer, v interface{}) error {
	var body *form.Values
	var commonParams *Params

	if params != nil {
		// This is a little unfortunate, but Go makes it impossible to compare
		// an interface value to nil without the use of the reflect package and
		// its true disciples insist that this is a feature and not a bug.
		//
		// Here we do invoke reflect because (1) we have to reflect anyway to
		// use encode with the form package, and (2) the corresponding removal
		// of boilerplate that this enables makes the small performance penalty
		// worth it.
		reflectValue := reflect.ValueOf(params)

		if reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil() {
			commonParams = params.GetParams()
			body = &form.Values{}
			form.AppendTo(body, params)
		}
	}

	return s.CallRaw(method, path, key, body, commonParams, v)
}

// CallMultipart is the Backend.CallMultipart implementation for invoking Stripe APIs.
func (s *BackendImplementation) CallMultipart(method, path, key, boundary string, body *bytes.Buffer, params *Params, v interface{}) error {
	contentType := "multipart/form-data; boundary=" + boundary

	req, err := s.NewRequest(method, path, key, contentType, params)
	if err != nil {
		return err
	}

	if err := s.Do(req, body, v); err != nil {
		return err
	}

	return nil
}

// CallRaw is the implementation for invoking Stripe APIs internally without a backend.
func (s *BackendImplementation) CallRaw(method, path, key string, form *form.Values, params *Params, v interface{}) error {
	var body string
	if form != nil && !form.Empty() {
		body = form.Encode()

		// On `GET`, move the payload into the URL
		if method == http.MethodGet {
			path += "?" + body
			body = ""
		}
	}
	bodyBuffer := bytes.NewBufferString(body)

	req, err := s.NewRequest(method, path, key, "application/x-www-form-urlencoded", params)
	if err != nil {
		return err
	}

	if err := s.Do(req, bodyBuffer, v); err != nil {
		return err
	}

	return nil
}

// NewRequest is used by Call to generate an http.Request. It handles encoding
// parameters and attaching the appropriate headers.
func (s *BackendImplementation) NewRequest(method, path, key, contentType string, params *Params) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = s.URL + path

	// Body is set later by `Do`.
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		if s.LogLevel > 0 {
			s.Logger.Printf("Cannot create Stripe request: %v\n", err)
		}
		return nil, err
	}

	authorization := "Bearer " + key

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Stripe-Version", apiversion)
	req.Header.Add("User-Agent", encodedUserAgent)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("X-Stripe-Client-User-Agent", encodedStripeUserAgent)

	if params != nil {
		if params.Context != nil {
			req = req.WithContext(params.Context)
		}

		if params.IdempotencyKey != nil {
			idempotencyKey := strings.TrimSpace(*params.IdempotencyKey)
			if len(idempotencyKey) > 255 {
				return nil, errors.New("cannot use an idempotency key longer than 255 characters")
			}

			req.Header.Add("Idempotency-Key", idempotencyKey)
		} else if isHTTPWriteMethod(method) {
			req.Header.Add("Idempotency-Key", NewIdempotencyKey())
		}

		if params.StripeAccount != nil {
			req.Header.Add("Stripe-Account", strings.TrimSpace(*params.StripeAccount))
		}

		for k, v := range params.Headers {
			for _, line := range v {
				req.Header.Add(k, line)
			}
		}
	}

	return req, nil
}

// Do is used by Call to execute an API request and parse the response. It uses
// the backend's HTTP client to execute the request and unmarshals the response
// into v. It also handles unmarshaling errors returned by the API.
func (s *BackendImplementation) Do(req *http.Request, body *bytes.Buffer, v interface{}) error {
	if s.LogLevel > 1 {
		s.Logger.Printf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)
	}

	var res *http.Response
	var err error
	for retry := 0; ; {
		start := time.Now()

		// This might look a little strange, but we set the request's body
		// outside of `NewRequest` so that we can get a fresh version every
		// time.
		//
		// The background is that back in the era of old style HTTP, it was
		// safe to reuse `Request` objects, but with the addition of HTTP/2,
		// it's now only sometimes safe. Reusing a `Request` with a body will
		// break.
		//
		// See some details here:
		//
		//     https://github.com/golang/go/issues/19653#issuecomment-341539160
		//
		// And our original bug report here:
		//
		//     https://github.com/stripe/stripe-go/issues/642
		//
		// To workaround the problem, we put a fresh `Body` onto the `Request`
		// every time we execute it, and this seems to empirically resolve the
		// problem.
		if body != nil {
			// We can safely reuse the same buffer that we used to encode our body,
			// but return a new reader to it everytime so that each read is from
			// the beginning.
			reader := bytes.NewReader(body.Bytes())

			req.Body = nopReadCloser{reader}

			// And also add the same thing to `Request.GetBody`, which allows
			// `net/http` to get a new body in cases like a redirect. This is
			// usually not used, but it doesn't hurt to set it in case it's
			// needed. See:
			//
			//     https://github.com/stripe/stripe-go/issues/710
			//
			req.GetBody = func() (io.ReadCloser, error) {
				reader := bytes.NewReader(body.Bytes())
				return nopReadCloser{reader}, nil
			}
		}

		res, err = s.HTTPClient.Do(req)

		if s.LogLevel > 2 {
			s.Logger.Printf("Request completed in %v (retry: %v)\n",
				time.Since(start), retry)
		}

		// If the response was okay, we're done, and it's safe to break out of
		// the retry loop.
		if !s.shouldRetry(err, res, retry) {
			break
		}

		if err != nil {
			if s.LogLevel > 0 {
				s.Logger.Printf("Request failed with error: %v\n", err)
			}
		} else {
			resBody, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				if s.LogLevel > 0 {
					s.Logger.Printf("Cannot read response: %v\n", err)
				}
			} else {
				if s.LogLevel > 0 {
					s.Logger.Printf("Request failed with body: %s (status: %v)\n", string(resBody), res.StatusCode)
				}
			}
		}

		sleepDuration := s.sleepTime(retry)
		retry++

		if s.LogLevel > 1 {
			s.Logger.Printf("Initiating retry %v for request %v %v%v after sleeping %v\n",
				retry, req.Method, req.URL.Host, req.URL.Path, sleepDuration)
		}

		time.Sleep(sleepDuration)
	}

	if err != nil {
		if s.LogLevel > 0 {
			s.Logger.Printf("Request failed: %v\n", err)
		}
		return err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if s.LogLevel > 0 {
			s.Logger.Printf("Cannot read response: %v\n", err)
		}
		return err
	}

	if res.StatusCode >= 400 {
		return s.ResponseToError(res, resBody)
	}

	if s.LogLevel > 2 {
		s.Logger.Printf("Response: %s\n", string(resBody))
	}

	if v != nil {
		return json.Unmarshal(resBody, v)
	}

	return nil
}

// ResponseToError converts a stripe response to an Error.
func (s *BackendImplementation) ResponseToError(res *http.Response, resBody []byte) error {
	var raw rawError
	if err := json.Unmarshal(resBody, &raw); err != nil {
		return err
	}
	// no error in resBody
	if raw.E == nil {
		err := errors.New(string(resBody))
		if s.LogLevel > 0 {
			s.Logger.Printf("Unparsable error returned from Stripe: %v\n", err)
		}
		return err
	}
	raw.E.HTTPStatusCode = res.StatusCode
	raw.E.RequestID = res.Header.Get("Request-Id")

	var typedError error
	switch raw.E.Type {
	case ErrorTypeAPI:
		typedError = &APIError{stripeErr: raw.E.Error}
	case ErrorTypeAPIConnection:
		typedError = &APIConnectionError{stripeErr: raw.E.Error}
	case ErrorTypeAuthentication:
		typedError = &AuthenticationError{stripeErr: raw.E.Error}
	case ErrorTypeCard:
		cardErr := &CardError{stripeErr: raw.E.Error}
		if raw.E.DeclineCode != nil {
			cardErr.DeclineCode = *raw.E.DeclineCode
		}
		typedError = cardErr
	case ErrorTypeInvalidRequest:
		typedError = &InvalidRequestError{stripeErr: raw.E.Error}
	case ErrorTypePermission:
		typedError = &PermissionError{stripeErr: raw.E.Error}
	case ErrorTypeRateLimit:
		typedError = &RateLimitError{stripeErr: raw.E.Error}
	}
	raw.E.Err = typedError

	if s.LogLevel > 0 {
		s.Logger.Printf("Error encountered from Stripe: %v\n", raw.E.Error)
	}
	return raw.E.Error
}

// SetMaxNetworkRetries sets max number of retries on failed requests
//
// This function is deprecated. Please use GetBackendWithConfig instead.
func (s *BackendImplementation) SetMaxNetworkRetries(maxNetworkRetries int) {
	s.MaxNetworkRetries = maxNetworkRetries
}

// SetNetworkRetriesSleep allows the normal sleep between network retries to be
// enabled or disabled.
//
// This function is available for internal testing only and should never be
// used in production.
func (s *BackendImplementation) SetNetworkRetriesSleep(sleep bool) {
	s.networkRetriesSleep = sleep
}

// Checks if an error is a problem that we should retry on. This includes both
// socket errors that may represent an intermittent problem and some special
// HTTP statuses.
func (s *BackendImplementation) shouldRetry(err error, resp *http.Response, numRetries int) bool {
	if numRetries >= s.MaxNetworkRetries {
		return false
	}

	if err != nil {
		return true
	}

	if resp.StatusCode == http.StatusConflict {
		return true
	}
	return false
}

// sleepTime calculates sleeping/delay time in milliseconds between failure and a new one request.
func (s *BackendImplementation) sleepTime(numRetries int) time.Duration {
	// We disable sleeping in some cases for tests.
	if !s.networkRetriesSleep {
		return 0 * time.Second
	}

	// Apply exponential backoff with minNetworkRetriesDelay on the
	// number of num_retries so far as inputs.
	delay := minNetworkRetriesDelay + minNetworkRetriesDelay*time.Duration(numRetries*numRetries)

	// Do not allow the number to exceed maxNetworkRetriesDelay.
	if delay > maxNetworkRetriesDelay {
		delay = maxNetworkRetriesDelay
	}

	// Apply some jitter by randomizing the value in the range of 75%-100%.
	jitter := rand.Int63n(int64(delay / 4))
	delay -= time.Duration(jitter)

	// But never sleep less than the base sleep seconds.
	if delay < minNetworkRetriesDelay {
		delay = minNetworkRetriesDelay
	}

	return delay
}

// Backends are the currently supported endpoints.
type Backends struct {
	API, Uploads Backend
	mu           sync.RWMutex
}

// Printfer is an interface to be implemented by Logger.
type Printfer interface {
	Printf(format string, v ...interface{})
}

// SupportedBackend is an enumeration of supported Stripe endpoints.
// Currently supported values are "api" and "uploads".
type SupportedBackend string

//
// Public functions
//

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// Float64 returns a pointer to the float64 value passed in.
func Float64(v float64) *float64 {
	return &v
}

// Float64Value returns the value of the float64 pointer passed in or
// 0 if the pointer is nil.
func Float64Value(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0
}

// FormatURLPath takes a format string (of the kind used in the fmt package)
// representing a URL path with a number of parameters that belong in the path
// and returns a formatted string.
//
// This is mostly a pass through to Sprintf. It exists to make it
// it impossible to accidentally provide a parameter type that would be
// formatted improperly; for example, a string pointer instead of a string.
//
// It also URL-escapes every given parameter. This usually isn't necessary for
// a standard Stripe ID, but is needed in places where user-provided IDs are
// allowed, like in coupons or plans. We apply it broadly for extra safety.
func FormatURLPath(format string, params ...string) string {
	// Convert parameters to interface{} and URL-escape them
	untypedParams := make([]interface{}, len(params))
	for i, param := range params {
		untypedParams[i] = interface{}(url.QueryEscape(param))
	}

	return fmt.Sprintf(format, untypedParams...)
}

// GetBackend returns one of the library's supported backends based off of the
// given argument.
//
// It returns an existing default backend if one's already been created.
func GetBackend(backendType SupportedBackend) Backend {
	var backend Backend

	backends.mu.RLock()
	switch backendType {
	case APIBackend:
		backend = backends.API
	case UploadsBackend:
		backend = backends.Uploads
	}
	backends.mu.RUnlock()
	if backend != nil {
		return backend
	}

	backend = GetBackendWithConfig(
		backendType,
		&BackendConfig{
			HTTPClient:        httpClient,
			LogLevel:          LogLevel,
			Logger:            Logger,
			MaxNetworkRetries: 0,
			URL:               "", // Set by GetBackendWithConfiguation when empty
		},
	)

	backends.mu.Lock()
	defer backends.mu.Unlock()

	switch backendType {
	case APIBackend:
		backends.API = backend
	case UploadsBackend:
		backends.Uploads = backend
	}

	return backend
}

// GetBackendWithConfig is the same as GetBackend except that it can be given a
// configuration struct that will configure certain aspects of the backend
// that's return.
func GetBackendWithConfig(backendType SupportedBackend, config *BackendConfig) Backend {
	if config.HTTPClient == nil {
		config.HTTPClient = httpClient
	}

	if config.Logger == nil {
		config.Logger = Logger
	}

	switch backendType {
	case APIBackend:
		if config.URL == "" {
			config.URL = apiURL
		}

		config.URL = normalizeURL(config.URL)

		return newBackendImplementation(backendType, config)

	case UploadsBackend:
		if config.URL == "" {
			config.URL = uploadsURL
		}

		config.URL = normalizeURL(config.URL)

		return newBackendImplementation(backendType, config)
	}

	return nil
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}

// Int64Value returns the value of the int64 pointer passed in or
// 0 if the pointer is nil.
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// NewBackends creates a new set of backends with the given HTTP client. You
// should only need to use this for testing purposes or on App Engine.
func NewBackends(httpClient *http.Client) *Backends {
	apiConfig := &BackendConfig{HTTPClient: httpClient}
	uploadConfig := &BackendConfig{HTTPClient: httpClient}
	return &Backends{
		API:     GetBackendWithConfig(APIBackend, apiConfig),
		Uploads: GetBackendWithConfig(UploadsBackend, uploadConfig),
	}
}

// ParseID attempts to parse a string scalar from a given JSON value which is
// still encoded as []byte. If the value was a string, it returns the string
// along with true as the second return value. If not, false is returned as the
// second return value.
//
// The purpose of this function is to detect whether a given value in a
// response from the Stripe API is a string ID or an expanded object.
func ParseID(data []byte) (string, bool) {
	s := string(data)

	if !strings.HasPrefix(s, "\"") {
		return "", false
	}

	if !strings.HasSuffix(s, "\"") {
		return "", false
	}

	return s[1 : len(s)-1], true
}

// SetAppInfo sets app information. See AppInfo.
func SetAppInfo(info *AppInfo) {
	if info != nil && info.Name == "" {
		panic(fmt.Errorf("App info name cannot be empty"))
	}
	appInfo = info

	// This is run in init, but we need to reinitialize it now that we have
	// some app info.
	initUserAgent()
}

// SetBackend sets the backend used in the binding.
func SetBackend(backend SupportedBackend, b Backend) {
	switch backend {
	case APIBackend:
		backends.API = b
	case UploadsBackend:
		backends.Uploads = b
	}
}

// SetHTTPClient overrides the default HTTP client.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.
func SetHTTPClient(client *http.Client) {
	httpClient = client
}

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

//
// Private constants
//

const apiURL = "https://api.stripe.com"

// apiversion is the currently supported API version
const apiversion = "2018-11-08"

// clientversion is the binding version
const clientversion = "55.3.0"

// defaultHTTPTimeout is the default timeout on the http.Client used by the library.
// This is chosen to be consistent with the other Stripe language libraries and
// to coordinate with other timeouts configured in the Stripe infrastructure.
const defaultHTTPTimeout = 80 * time.Second

// maxNetworkRetriesDelay and minNetworkRetriesDelay defines sleep time in milliseconds between
// tries to send HTTP request again after network failure.
const maxNetworkRetriesDelay = 5000 * time.Millisecond
const minNetworkRetriesDelay = 500 * time.Millisecond

const uploadsURL = "https://uploads.stripe.com"

//
// Private types
//

// nopReadCloser's sole purpose is to give us a way to turn an `io.Reader` into
// an `io.ReadCloser` by adding a no-op implementation of the `Closer`
// interface. We need this because `http.Request`'s `Body` takes an
// `io.ReadCloser` instead of a `io.Reader`.
type nopReadCloser struct {
	io.Reader
}

func (nopReadCloser) Close() error { return nil }

// stripeClientUserAgent contains information about the current runtime which
// is serialized and sent in the `X-Stripe-Client-User-Agent` as additional
// debugging information.
type stripeClientUserAgent struct {
	Application     *AppInfo `json:"application"`
	BindingsVersion string   `json:"bindings_version"`
	Language        string   `json:"lang"`
	LanguageVersion string   `json:"lang_version"`
	Publisher       string   `json:"publisher"`
	Uname           string   `json:"uname"`
}

//
// Private variables
//

var appInfo *AppInfo
var backends Backends
var encodedStripeUserAgent string
var encodedUserAgent string
var httpClient = &http.Client{Timeout: defaultHTTPTimeout}

//
// Private functions
//

// getUname tries to get a uname from the system, but not that hard. It tries
// to execute `uname -a`, but swallows any errors in case that didn't work
// (i.e. non-Unix non-Mac system or some other reason).
func getUname() string {
	path, err := exec.LookPath("uname")
	if err != nil {
		return UnknownPlatform
	}

	cmd := exec.Command(path, "-a")
	var out bytes.Buffer
	cmd.Stderr = nil // goes to os.DevNull
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return UnknownPlatform
	}

	return out.String()
}

func init() {
	Logger = log.New(os.Stderr, "", log.LstdFlags)
	initUserAgent()
}

func initUserAgent() {
	encodedUserAgent = "Stripe/v1 GoBindings/" + clientversion
	if appInfo != nil {
		encodedUserAgent += " " + appInfo.formatUserAgent()
	}

	stripeUserAgent := &stripeClientUserAgent{
		Application:     appInfo,
		BindingsVersion: clientversion,
		Language:        "go",
		LanguageVersion: runtime.Version(),
		Publisher:       "stripe",
		Uname:           getUname(),
	}
	marshaled, err := json.Marshal(stripeUserAgent)
	// Encoding this struct should never be a problem, so we're okay to panic
	// in case it is for some reason.
	if err != nil {
		panic(err)
	}
	encodedStripeUserAgent = string(marshaled)
}

func isHTTPWriteMethod(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}

// newBackendImplementation returns a new Backend based off a given type and
// fully initialized BackendConfig struct.
//
// The vast majority of the time you should be calling GetBackendWithConfig
// instead of this function.
func newBackendImplementation(backendType SupportedBackend, config *BackendConfig) Backend {
	return &BackendImplementation{
		HTTPClient:          config.HTTPClient,
		LogLevel:            config.LogLevel,
		Logger:              config.Logger,
		MaxNetworkRetries:   config.MaxNetworkRetries,
		Type:                backendType,
		URL:                 config.URL,
		networkRetriesSleep: true,
	}
}

func normalizeURL(url string) string {
	// All paths include a leading slash, so to keep logs pretty, trim a
	// trailing slash on the URL.
	url = strings.TrimSuffix(url, "/")

	// For a long time we had the `/v1` suffix as part of a configured URL
	// rather than in the per-package URLs throughout the library. Continue
	// to support this for the time being by stripping one that's been
	// passed for better backwards compatibility.
	url = strings.TrimSuffix(url, "/v1")

	return url
}
