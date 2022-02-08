package acrhelper

import (
	"context"
	"fmt"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/services/containerregistry"
)

type ACRHelper struct {
	opts       *Options
	authorizer auth.Authorizer
}

var _ credentials.Helper = (*ACRHelper)(nil)

// NewACRHelper will return a credentials.Helper that has implemented Get() for Azure Container Registry
func NewACRHelper(setters ...Option) credentials.Helper {
	opts := newOptions(setters...)

	var authorizer auth.Authorizer
	if opts.cacheAuthorizerAtConstruction {
		ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
		defer cancel()

		a, err := newAuthorizer(ctx, opts)
		if err == nil {
			authorizer = a
		}
	}

	return &ACRHelper{
		opts,
		authorizer,
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
	ctx, cancel := context.WithTimeout(context.Background(), acr.opts.timeout)
	defer cancel()

	authorizer := acr.authorizer
	if authorizer == nil {
		a, err := newAuthorizer(ctx, acr.opts)
		if err != nil {
			return "", "", err
		}

		authorizer = a
	}

	cr := containerregistry.NewContainerRegistryClient(authorizer, serverURL, acr.opts.tenantID)
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

// newAuthorizer will try to create a cached authorizer based on the config provided
// 1. Client Certificate
// 2. Client Secret
// 3. MSI
// 4. GitHub OIDC
// 5. Azure CLI
func newAuthorizer(ctx context.Context, opts *Options) (auth.Authorizer, error) {
	authorizer, err := auth.NewClientCertificateAuthorizer(ctx, *opts.environment, opts.environment.ResourceManager, opts.tokenVersion, opts.tenantID, nil, opts.clientID, nil, opts.clientCertPath, opts.clientCertPassword)
	if err == nil {
		_, err := authorizer.Token()
		if err == nil {
			return authorizer, nil
		}
	}

	if opts.tenantID != "" && opts.clientID != "" && opts.clientSecret != "" {
		authorizer, err = auth.NewClientSecretAuthorizer(ctx, *opts.environment, opts.environment.ResourceManager, opts.tokenVersion, opts.tenantID, nil, opts.clientID, opts.clientSecret)
		if err == nil {
			_, err := authorizer.Token()
			if err == nil {
				return authorizer, nil
			}
		}
	}

	authorizer, err = auth.NewMsiAuthorizer(ctx, opts.environment.ResourceManager, "", opts.clientID)
	if err == nil {
		_, err := authorizer.Token()
		if err == nil {
			return authorizer, nil
		}
	}

	if opts.gitHubRequestToken != "" || opts.gitHubRequestTokenURL != "" {
		authorizer, err = auth.NewGitHubOIDCAuthorizer(ctx, *opts.environment, opts.environment.ResourceManager, opts.tenantID, nil, opts.clientID, opts.gitHubRequestTokenURL, opts.gitHubRequestToken)
		if err == nil {
			_, err := authorizer.Token()
			if err == nil {
				return authorizer, nil
			}
		}
	}

	authorizer, err = auth.NewAzureCliAuthorizer(ctx, opts.environment.ResourceManager, opts.tenantID)
	if err == nil {
		_, err := authorizer.Token()
		if err == nil {
			return authorizer, nil
		}
	}

	return nil, fmt.Errorf("no valid authorizer could be found")
}
