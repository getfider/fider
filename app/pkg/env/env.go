package env

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"path"

	"github.com/getfider/fider/app/pkg/errors"
	"github.com/joeshaw/envdecode"
)

var (
	// this value is set by the CI build
	commithash = ""
	version    = "0.22.0"
)

func Version() string {
	if commithash == "" {
		return fmt.Sprintf("%s-dev", version)
	}

	return fmt.Sprintf("%s-%s", version, commithash)
}

type config struct {
	Environment    string `env:"GO_ENV,default=production"`
	SignUpDisabled bool   `env:"SIGNUP_DISABLED,default=false"`
	TLS            struct {
		Automatic      bool   `env:"SSL_AUTO,default=false"`
		Certificate    string `env:"SSL_CERT"`
		CertificateKey string `env:"SSL_CERT_KEY"`
	}
	HTTP struct {
		ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,default=5s,strict"`
		WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,default=10s,strict"`
		IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,default=120s,strict"`
	}
	Port       string `env:"PORT,default=3000"`
	HostMode   string `env:"HOST_MODE,default=single"`
	HostDomain string `env:"HOST_DOMAIN"`
	BaseURL    string `env:"BASE_URL"`
	Locale     string `env:"LOCALE,default=en"`
	JWTSecret  string `env:"JWT_SECRET,required"`
	Paddle     struct {
		IsSandbox      bool   `env:"PADDLE_SANDBOX,default=false"`
		VendorID       string `env:"PADDLE_VENDOR_ID"`
		VendorAuthCode string `env:"PADDLE_VENDOR_AUTHCODE"`
		PlanID         string `env:"PADDLE_PLAN_ID"`
		PublicKey      string `env:"PADDLE_PUBLIC_KEY"`
	}
	Metrics struct {
		Enabled bool   `env:"METRICS_ENABLED,default=false"`
		Port    string `env:"METRICS_PORT,default=4000"`
	}
	Database struct {
		URL          string `env:"DATABASE_URL,required"`
		MaxIdleConns int    `env:"DATABASE_MAX_IDLE_CONNS,default=2,strict"`
		MaxOpenConns int    `env:"DATABASE_MAX_OPEN_CONNS,default=4,strict"`
	}
	CDN struct {
		Host string `env:"CDN_HOST"`
	}
	Log struct {
		Level      string `env:"LOG_LEVEL,default=INFO"`
		Structured bool   `env:"LOG_STRUCTURED,default=false"`
		Console    bool   `env:"LOG_CONSOLE,default=true"`
		Sql        bool   `env:"LOG_SQL,default=true"`
		File       bool   `env:"LOG_FILE,default=false"`
		OutputFile string `env:"LOG_FILE_OUTPUT,default=logs/output.log"`
	}
	OAuth struct {
		Google struct {
			ClientID string `env:"OAUTH_GOOGLE_CLIENTID"`
			Secret   string `env:"OAUTH_GOOGLE_SECRET"`
		}
		Facebook struct {
			AppID  string `env:"OAUTH_FACEBOOK_APPID"`
			Secret string `env:"OAUTH_FACEBOOK_SECRET"`
		}
		GitHub struct {
			ClientID string `env:"OAUTH_GITHUB_CLIENTID"`
			Secret   string `env:"OAUTH_GITHUB_SECRET"`
		}
	}
	Email struct {
		Type      string `env:"EMAIL"` // possible values: smtp, mailgun, awsses
		NoReply   string `env:"EMAIL_NOREPLY,required"`
		Allowlist string `env:"EMAIL_ALLOWLIST"`
		Blocklist string `env:"EMAIL_BLOCKLIST"`
		AWSSES    struct {
			Region          string `env:"EMAIL_AWSSES_REGION"`
			AccessKeyID     string `env:"EMAIL_AWSSES_ACCESS_KEY_ID"`
			SecretAccessKey string `env:"EMAIL_AWSSES_SECRET_ACCESS_KEY"`
		}
		Mailgun struct {
			APIKey string `env:"EMAIL_MAILGUN_API"`
			Domain string `env:"EMAIL_MAILGUN_DOMAIN"`
			Region string `env:"EMAIL_MAILGUN_REGION,default=US"` // possible values: US or EU
		}
		SMTP struct {
			Host           string `env:"EMAIL_SMTP_HOST"`
			Port           string `env:"EMAIL_SMTP_PORT"`
			Username       string `env:"EMAIL_SMTP_USERNAME"`
			Password       string `env:"EMAIL_SMTP_PASSWORD"`
			EnableStartTLS bool   `env:"EMAIL_SMTP_ENABLE_STARTTLS,default=true"`
		}
	}
	BlobStorage struct {
		Type string `env:"BLOB_STORAGE,default=sql"` // possible values: sql, fs or s3
		S3   struct {
			EndpointURL     string `env:"BLOB_STORAGE_S3_ENDPOINT_URL"`
			Region          string `env:"BLOB_STORAGE_S3_REGION"`
			AccessKeyID     string `env:"BLOB_STORAGE_S3_ACCESS_KEY_ID"`
			SecretAccessKey string `env:"BLOB_STORAGE_S3_SECRET_ACCESS_KEY"`
			BucketName      string `env:"BLOB_STORAGE_S3_BUCKET"`
		}
		FS struct {
			Path string `env:"BLOB_STORAGE_FS_PATH"`
		}
	}
	Maintenance struct {
		Enabled bool   `env:"MAINTENANCE,default=false,strict"`
		Message string `env:"MAINTENANCE_MESSAGE"`
		Until   string `env:"MAINTENANCE_UNTIL"`
	}
	GoogleAnalytics string `env:"GOOGLE_ANALYTICS"`
}

// Config is a strongly typed reference to all configuration parsed from Environment Variables
var Config config

func init() {
	Reload()
}

// Reload configuration from current Enviornment Variables
func Reload() {
	Config = config{}
	err := envdecode.Decode(&Config)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse environment variables"))
	}

	if IsSingleHostMode() {
		if Config.HostDomain != "" {
			panic("HOST_DOMAIN environment variable has been replaced by BASE_URL. Set it to your site base url, e.g: https://feedback.mysite.com")
		}

		mustBeSet("BASE_URL")

		_, err := url.Parse(Config.BaseURL)
		if err != nil {
			panic(errors.Wrap(err, "'%s' is not a valid URL", Config.BaseURL))
		}
	}

	// Email Type can be inferred if absense
	if Config.Email.Type == "" {
		if Config.Email.Mailgun.APIKey != "" {
			Config.Email.Type = "mailgun"
		} else if Config.Email.AWSSES.AccessKeyID != "" {
			Config.Email.Type = "awsses"
		} else {
			Config.Email.Type = "smtp"
		}
	}

	emailType := Config.Email.Type
	if emailType == "mailgun" {
		mustBeSet("EMAIL_MAILGUN_API")
		mustBeSet("EMAIL_MAILGUN_DOMAIN")
	} else if emailType == "awsses" {
		mustBeSet("EMAIL_AWSSES_REGION")
		mustBeSet("EMAIL_AWSSES_ACCESS_KEY_ID")
		mustBeSet("EMAIL_AWSSES_SECRET_ACCESS_KEY")
	} else if emailType == "smtp" {
		mustBeSet("EMAIL_SMTP_HOST")
		mustBeSet("EMAIL_SMTP_PORT")
	}

	bsType := Config.BlobStorage.Type
	if bsType == "s3" {
		mustBeSet("BLOB_STORAGE_S3_BUCKET")
	} else if bsType == "fs" {
		mustBeSet("BLOB_STORAGE_FS_PATH")
	}
}

func mustBeSet(name string) {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Errorf("could not find environment variable named '%s'", name))
	}
}

// IsSingleHostMode returns true if host mode is set to single tenant
func IsSingleHostMode() bool {
	return Config.HostMode == "single"
}

var hasLegal *bool

// HasLegal returns true if current instance contains legal documents: privacy.md and terms.md
func HasLegal() bool {
	if hasLegal == nil {
		_, err1 := os.Stat(Etc("privacy.md"))
		_, err2 := os.Stat(Etc("terms.md"))
		exists := err1 == nil && err2 == nil
		hasLegal = &exists
	}
	return *hasLegal
}

// MultiTenantDomain returns domain name of current instance for multi tenant hosts
func MultiTenantDomain() string {
	if !IsSingleHostMode() {
		return "." + Config.HostDomain
	}
	return ""
}

// IsBillingEnabled returns true if Paddle is configured
func IsBillingEnabled() bool {
	return Config.Paddle.VendorID != "" && Config.Paddle.VendorAuthCode != ""
}

// IsProduction returns true on Fider production environment
func IsProduction() bool {
	return Config.Environment == "production" || (!IsTest() && !IsDevelopment())
}

// IsTest returns true on Fider test environment
func IsTest() bool {
	return Config.Environment == "test"
}

// IsDevelopment returns true on Fider production environment
func IsDevelopment() bool {
	return Config.Environment == "development"
}

// Path returns root path of project + given path
func Path(p ...string) string {
	root := "./"
	if IsTest() {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		root = path.Join(basepath, "../../../")
	}

	elems := append([]string{root}, p...)
	return path.Join(elems...)
}

// Etc returns a path to a folder or file inside the /etc/ folder
func Etc(p ...string) string {
	paths := append([]string{"etc"}, p...)
	return Path(paths...)
}

// Subdomain returns the Fider subdomain (if available) from given host
func Subdomain(host string) string {
	if IsSingleHostMode() {
		return ""
	}

	domain := MultiTenantDomain()
	if domain != "" && strings.Contains(host, domain) {
		return strings.Replace(host, domain, "", -1)
	}

	if Config.CDN.Host != "" {
		domain = Config.CDN.Host
		parts := strings.Split(domain, ":")
		if parts[0] != "" && strings.Contains(host, "."+parts[0]) {
			return strings.Replace(host, "."+parts[0], "", -1)
		}
	}

	return ""
}
