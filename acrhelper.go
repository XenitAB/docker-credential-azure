package acrhelper

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/services/containerregistry"
)

type config struct {
	*auth.Config
	gitHubTokenURL string
	gitHubToken    string
	timeout        time.Duration
}

type ACRHelper struct {
	cfg              config
	cachedAuthorizer auth.Authorizer
}

var _ credentials.Helper = (*ACRHelper)(nil)

// NewACRHelper will return a credentials.Helper that has implemented Get() for Azure Container Registry
// Parameters:
// timeout = what the timeout will be for the contexts created when using this helper
// cacheAuthorizer = should the constructor try to create a cached authorizer
func NewACRHelper(timeout time.Duration, cacheAuthorizer bool) credentials.Helper {
	cfg := newConfig(timeout)

	var cachedAuthorizer auth.Authorizer
	if cacheAuthorizer {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		authorizer, err := newCachedAuthorizer(ctx, cfg)
		if err == nil {
			cachedAuthorizer = authorizer
		}
	}

	return &ACRHelper{
		cfg,
		cachedAuthorizer,
	}
}

// Add appends credentials to the store.
// NOT implemented.
func (ACRHelper) Add(_ *credentials.Credentials) error {
	return fmt.Errorf("method Add() is not implemented")
}

// Delete removes credentials from the store.
// NOT implemented.
func (ACRHelper) Delete(_ string) error {
	return fmt.Errorf("method Delete() not implemented")
}

// Get retrieves credentials from the store.
// It returns username and secret as strings.
func (acr ACRHelper) Get(serverURL string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), acr.cfg.timeout)
	defer cancel()

	cachedAuthorizer := acr.cachedAuthorizer
	if cachedAuthorizer == nil {
		a, err := newCachedAuthorizer(ctx, acr.cfg)
		if err != nil {
			return "", "", err
		}

		cachedAuthorizer = a
	}

	cr := containerregistry.NewContainerRegistryClient(cachedAuthorizer, serverURL, acr.cfg.TenantID)
	acrToken, _, err := cr.ExchangeRefreshToken(ctx)
	if err != nil {
		return "", "", err
	}

	return "<token>", acrToken, nil
}

// List returns the stored serverURLs and their associated usernames.
// NOT implemented.
func (acr ACRHelper) List() (map[string]string, error) {
	return nil, fmt.Errorf("method List() is not implemented")
}

// newCachedAuthorizer will try to create a cached authorizer based on the config provided
// 1. Client Certificate
// 2. Client Secret
// 3. MSI
// 4. GitHub OIDC
// 5. Azure CLI
func newCachedAuthorizer(ctx context.Context, cfg config) (auth.Authorizer, error) {
	authorizer, err := auth.NewClientCertificateAuthorizer(ctx, cfg.Environment, cfg.Environment.ResourceManager, cfg.Version, cfg.TenantID, cfg.AuxiliaryTenantIDs, cfg.ClientID, cfg.ClientCertData, cfg.ClientCertPath, cfg.ClientCertPassword)
	if err == nil {
		cachedAuthorizer := auth.NewCachedAuthorizer(authorizer)
		_, err := cachedAuthorizer.Token()
		if err == nil {
			return cachedAuthorizer, nil
		}
	}

	if cfg.TenantID != "" && cfg.ClientID != "" && cfg.ClientSecret != "" {
		authorizer, err = auth.NewClientSecretAuthorizer(ctx, cfg.Environment, cfg.Environment.ResourceManager, cfg.Version, cfg.TenantID, cfg.AuxiliaryTenantIDs, cfg.ClientID, cfg.ClientSecret)
		if err == nil {
			cachedAuthorizer := auth.NewCachedAuthorizer(authorizer)
			_, err := cachedAuthorizer.Token()
			if err == nil {
				return cachedAuthorizer, nil
			}
		}
	}

	authorizer, err = auth.NewMsiAuthorizer(ctx, cfg.Environment.ResourceManager, cfg.MsiEndpoint, cfg.ClientID)
	if err == nil {
		cachedAuthorizer := auth.NewCachedAuthorizer(authorizer)
		_, err := cachedAuthorizer.Token()
		if err == nil {
			return cachedAuthorizer, nil
		}
	}

	if cfg.gitHubToken != "" || cfg.gitHubTokenURL != "" {
		authorizer, err = auth.NewGitHubOIDCAuthorizer(ctx, cfg.Environment, cfg.Environment.ResourceManager, cfg.TenantID, cfg.AuxiliaryTenantIDs, cfg.ClientID, cfg.gitHubTokenURL, cfg.gitHubToken)
		if err == nil {
			cachedAuthorizer := auth.NewCachedAuthorizer(authorizer)
			_, err := cachedAuthorizer.Token()
			if err == nil {
				return cachedAuthorizer, nil
			}
		}
	}

	authorizer, err = auth.NewAzureCliAuthorizer(ctx, cfg.Environment.ResourceManager, cfg.TenantID)
	if err == nil {
		cachedAuthorizer := auth.NewCachedAuthorizer(authorizer)
		_, err := cachedAuthorizer.Token()
		if err == nil {
			return cachedAuthorizer, nil
		}
	}

	return nil, fmt.Errorf("no valid authorizer could be found")
}

// newConfig will create a config based on the environment
// Environment variables read:
// AZURE_TENANT_ID = Azure tenant ID
// AZURE_CLIENT_ID = Client ID for the Azure AD application
// AZURE_CLIENT_SECRET = Client secret for the Azure AD application
// AZURE_CERTIFICATE_PATH = Path to the certificate
// AZURE_CERTIFICATE_PASSWORD = Password to the certificate
// ACTIONS_ID_TOKEN_REQUEST_URL = Provided by GitHub Actions
// ACTIONS_ID_TOKEN_REQUEST_TOKEN = Provided by GitHub Actions
func newConfig(timeout time.Duration) config {
	authCfg := &auth.Config{
		Environment:        environments.Global,
		TenantID:           strings.TrimSpace(os.Getenv("AZURE_TENANT_ID")),
		ClientID:           strings.TrimSpace(os.Getenv("AZURE_CLIENT_ID")),
		ClientSecret:       strings.TrimSpace(os.Getenv("AZURE_CLIENT_SECRET")),
		ClientCertPath:     os.Getenv("AZURE_CERTIFICATE_PATH"),
		ClientCertPassword: os.Getenv("AZURE_CERTIFICATE_PASSWORD"),
		Version:            auth.TokenVersion2,
	}

	gitHubTokenURL := strings.TrimSpace(os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"))
	gitHubToken := strings.TrimSpace(os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN"))

	return config{
		authCfg,
		gitHubTokenURL,
		gitHubToken,
		timeout,
	}
}
