package acrhelper

import (
	"os"
	"strings"
	"time"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

// Options for the ACRHelper
type Options struct {
	environment                   *environments.Environment
	tokenVersion                  auth.TokenVersion
	tenantID                      string
	clientID                      string
	clientSecret                  string
	clientCertPath                string
	clientCertPassword            string
	gitHubRequestTokenURL         string
	gitHubRequestToken            string
	timeout                       time.Duration
	cacheAuthorizerAtConstruction bool
}

// Option is a function to set Options
type Option func(*Options)

func newOptions(setters ...Option) *Options {
	opts := &Options{
		environment:                   &environments.Global,
		tokenVersion:                  auth.TokenVersion2,
		tenantID:                      strings.TrimSpace(os.Getenv("AZURE_TENANT_ID")),
		clientID:                      strings.TrimSpace(os.Getenv("AZURE_CLIENT_ID")),
		clientSecret:                  strings.TrimSpace(os.Getenv("AZURE_CLIENT_SECRET")),
		clientCertPath:                os.Getenv("AZURE_CERTIFICATE_PATH"),
		clientCertPassword:            os.Getenv("AZURE_CERTIFICATE_PASSWORD"),
		gitHubRequestTokenURL:         strings.TrimSpace(os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")),
		gitHubRequestToken:            strings.TrimSpace(os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")),
		timeout:                       10 * time.Second,
		cacheAuthorizerAtConstruction: true,
	}

	for _, setter := range setters {
		setter(opts)
	}

	return opts
}

// WithEnvironment will set the Azure environment to use
// Default: &environments.Global
func WithEnvironment(opt *environments.Environment) Option {
	return func(opts *Options) {
		opts.environment = opt
	}
}

// WithTokenVersion will set the Azure AD token version to request
// Default: auth.TokenVersion2
func WithTokenVersion(opt auth.TokenVersion) Option {
	return func(opts *Options) {
		opts.tokenVersion = opt
	}
}

// WithTenantID will set the tenant id to use for requesting Azure AD tokens
// Default (environment variable): AZURE_TENANT_ID
func WithTenantID(opt string) Option {
	return func(opts *Options) {
		opts.tenantID = opt
	}
}

// WithClientID will set the client id to use for requesting Azure AD tokens
// Default (environment variable): AZURE_CLIENT_ID
func WithClientID(opt string) Option {
	return func(opts *Options) {
		opts.clientID = opt
	}
}

// WithClientSecret will set the client secret to use for requesting Azure AD tokens
// Default (environment variable): AZURE_CLIENT_SECRET
func WithClientSecret(opt string) Option {
	return func(opts *Options) {
		opts.clientSecret = opt
	}
}

// WithClientCertPath will set the client certificate path to use for requesting Azure AD tokens
// Default (environment variable): AZURE_CERTIFICATE_PATH
func WithClientCertPath(opt string) Option {
	return func(opts *Options) {
		opts.clientCertPath = opt
	}
}

// WithClientCertPassword will set the client certificate password to use for requesting Azure AD tokens
// Default (environment variable): AZURE_CERTIFICATE_PASSWORD
func WithClientCertPassword(opt string) Option {
	return func(opts *Options) {
		opts.clientCertPassword = opt
	}
}

// WithGitHubRequestTokenURL will set the GitHub request token URL to use for requesting Azure AD tokens
// Default (environment variable): ACTIONS_ID_TOKEN_REQUEST_URL
func WithGitHubRequestTokenURL(opt string) Option {
	return func(opts *Options) {
		opts.gitHubRequestTokenURL = opt
	}
}

// WithGitHubRequestToken will set the GitHub request token to use for requesting Azure AD tokens
// Default (environment variable): ACTIONS_ID_TOKEN_REQUEST_TOKEN
func WithGitHubRequestToken(opt string) Option {
	return func(opts *Options) {
		opts.gitHubRequestToken = opt
	}
}

// WithTimeout will set the timeout for contexts created inside of ACRHelper
// Default: 10 seconds
func WithTimeout(opt time.Duration) Option {
	return func(opts *Options) {
		opts.timeout = opt
	}
}

// WithCacheAuthorizerAtConstruction will, if enabled, try to create an Azure AD authorizer at construction (`NewACRHelper()`)
// Default: true
func WithCacheAuthorizerAtConstruction(opt bool) Option {
	return func(opts *Options) {
		opts.cacheAuthorizerAtConstruction = opt
	}
}
