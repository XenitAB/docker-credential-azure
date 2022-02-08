package acrhelper

import (
	"testing"
	"time"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

func TestOptions(t *testing.T) {
	expectedResult := &Options{
		environment:                   &environments.Canary,
		tokenVersion:                  auth.TokenVersion1,
		tenantID:                      "ze-tenant-id",
		clientID:                      "ze-client-id",
		clientSecret:                  "ze-client-secret",
		clientCertPath:                "ze-client-cert-path",
		clientCertPassword:            "ze-client-cert-password",
		gitHubRequestTokenURL:         "ze-github-request-token-url",
		gitHubRequestToken:            "ze-github-request-token",
		timeout:                       1337 * time.Millisecond,
		cacheAuthorizerAtConstruction: true,
	}

	setters := []Option{
		WithEnvironment(&environments.Canary),
		WithTokenVersion(auth.TokenVersion1),
		WithTenantID("ze-tenant-id"),
		WithClientID("ze-client-id"),
		WithClientSecret("ze-client-secret"),
		WithClientCertPath("ze-client-cert-path"),
		WithClientCertPassword("ze-client-cert-password"),
		WithGitHubRequestTokenURL("ze-github-request-token-url"),
		WithGitHubRequestToken("ze-github-request-token"),
		WithTimeout(1337 * time.Millisecond),
		WithCacheAuthorizerAtConstruction(true),
	}

	result := &Options{}

	for _, setter := range setters {
		setter(result)
	}

	if expectedResult.environment != result.environment {
		t.Fatalf("want %v got %v", expectedResult.environment, result.environment)
	}

	if expectedResult.tokenVersion != result.tokenVersion {
		t.Fatalf("want %d got %d", expectedResult.tokenVersion, result.tokenVersion)
	}

	if expectedResult.tenantID != result.tenantID {
		t.Fatalf("want %q got %q", expectedResult.tenantID, result.tenantID)
	}

	if expectedResult.clientID != result.clientID {
		t.Fatalf("want %q got %q", expectedResult.clientID, result.clientID)
	}

	if expectedResult.clientSecret != result.clientSecret {
		t.Fatalf("want %q got %q", expectedResult.clientSecret, result.clientSecret)
	}

	if expectedResult.clientCertPath != result.clientCertPath {
		t.Fatalf("want %q got %q", expectedResult.clientCertPath, result.clientCertPath)
	}

	if expectedResult.clientCertPassword != result.clientCertPassword {
		t.Fatalf("want %q got %q", expectedResult.clientCertPassword, result.clientCertPassword)
	}

	if expectedResult.gitHubRequestTokenURL != result.gitHubRequestTokenURL {
		t.Fatalf("want %q got %q", expectedResult.gitHubRequestTokenURL, result.gitHubRequestTokenURL)
	}

	if expectedResult.gitHubRequestToken != result.gitHubRequestToken {
		t.Fatalf("want %q got %q", expectedResult.gitHubRequestToken, result.gitHubRequestToken)
	}

	if expectedResult.timeout != result.timeout {
		t.Fatalf("want %q got %q", expectedResult.timeout.String(), result.timeout.String())
	}

	if expectedResult.cacheAuthorizerAtConstruction != result.cacheAuthorizerAtConstruction {
		t.Fatalf("want %t got %t", expectedResult.cacheAuthorizerAtConstruction, result.cacheAuthorizerAtConstruction)
	}
}
