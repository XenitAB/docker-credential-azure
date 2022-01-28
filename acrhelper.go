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

type ACRHelper struct {
	cfg            *auth.Config
	gitHubTokenURL string
	gitHubToken    string
	timeout        time.Duration
}

var _ credentials.Helper = (*ACRHelper)(nil)

func NewACRHelper() credentials.Helper {
	authConfig := newConfig()

	return &ACRHelper{
		cfg:            authConfig,
		gitHubTokenURL: strings.TrimSpace(os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")),
		gitHubToken:    strings.TrimSpace(os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")),
		timeout:        30 * time.Second,
	}
}

func (ACRHelper) Add(_ *credentials.Credentials) error {
	return fmt.Errorf("method Add() is not implemented")
}

func (ACRHelper) Delete(_ string) error {
	return fmt.Errorf("method Delete() not implemented")
}

func (acr ACRHelper) Get(serverURL string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), acr.timeout)
	defer cancel()

	authorizer, err := acr.newAuthorizer(ctx)
	if err != nil {
		return "", "", err
	}

	cr := containerregistry.NewContainerRegistryClient(authorizer, serverURL, acr.cfg.TenantID)
	acrToken, err := cr.ExchangeToken(ctx)
	if err != nil {
		return "", "", err
	}

	return "<token>", acrToken, nil
}

func (acr ACRHelper) List() (map[string]string, error) {
	return nil, fmt.Errorf("method List() is not implemented")
}

func (acr ACRHelper) newAuthorizer(ctx context.Context) (auth.Authorizer, error) {
	authCfg := acr.cfg
	a, err := auth.NewClientCertificateAuthorizer(ctx, authCfg.Environment, authCfg.Environment.ResourceManager, authCfg.Version, authCfg.TenantID, authCfg.AuxiliaryTenantIDs, authCfg.ClientID, authCfg.ClientCertData, authCfg.ClientCertPath, authCfg.ClientCertPassword)
	if err == nil {
		return a, nil
	}

	if authCfg.TenantID != "" && authCfg.ClientID != "" && authCfg.ClientSecret != "" {
		a, err = auth.NewClientSecretAuthorizer(ctx, authCfg.Environment, authCfg.Environment.ResourceManager, authCfg.Version, authCfg.TenantID, authCfg.AuxiliaryTenantIDs, authCfg.ClientID, authCfg.ClientSecret)
		if err == nil {
			return a, nil
		}
	}

	a, err = auth.NewMsiAuthorizer(ctx, authCfg.Environment.ResourceManager, authCfg.MsiEndpoint, authCfg.ClientID)
	if err == nil {
		return a, nil
	}

	if acr.gitHubToken != "" || acr.gitHubTokenURL != "" {
		a, err = auth.NewGitHubOIDCAuthorizer(ctx, authCfg.Environment, authCfg.Environment.ResourceManager, authCfg.TenantID, authCfg.AuxiliaryTenantIDs, authCfg.ClientID, acr.gitHubTokenURL, acr.gitHubToken)
		if err == nil {
			return a, nil
		}
	}

	a, err = auth.NewAzureCliAuthorizer(ctx, authCfg.Environment.ResourceManager, authCfg.TenantID)
	if err == nil {
		return a, nil
	}

	return nil, fmt.Errorf("no valid authorizer could be found")
}

func newConfig() *auth.Config {
	auxTenants := strings.Split(os.Getenv("AZURE_AUXILIARY_TENANT_IDS"), ";")
	for i := range auxTenants {
		auxTenants[i] = strings.TrimSpace(auxTenants[i])
	}

	return &auth.Config{
		Environment:        environments.Global,
		TenantID:           strings.TrimSpace(os.Getenv("AZURE_TENANT_ID")),
		ClientID:           strings.TrimSpace(os.Getenv("AZURE_CLIENT_ID")),
		ClientSecret:       strings.TrimSpace(os.Getenv("AZURE_CLIENT_SECRET")),
		ClientCertPath:     os.Getenv("AZURE_CERTIFICATE_PATH"),
		ClientCertPassword: os.Getenv("AZURE_CERTIFICATE_PASSWORD"),
		AuxiliaryTenantIDs: auxTenants,
		Version:            auth.TokenVersion2,
	}
}
