package acrhelper

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/manicminer/hamilton/services/containerregistry"
)

func TestCI(t *testing.T) {
	if strings.EqualFold(os.Getenv("CI"), "true") {
		if os.Getenv("CONTAINER_REGISTRY_NAME") == "" {
			t.Fatal("when running in CI, the environment variable 'CONTAINER_REGISTRY_NAME' needs to be set")
		}
	}
}

func TestACRHelperE2E(t *testing.T) {
	containerRegistryName := os.Getenv("CONTAINER_REGISTRY_NAME")
	if containerRegistryName == "" {
		t.Skip("environment variable CONTAINER_REGISTRY_NAME not set")
	}

	serverURL := fmt.Sprintf("%s.azurecr.io", containerRegistryName)

	// without caching the authorizer at construction
	_, _ = testACRHelperGet(t, serverURL, WithTimeout(30*time.Second), WithCacheAuthorizerAtConstruction(false))

	// with caching the authorizer at construction
	helper, password := testACRHelperGet(t, serverURL, WithTimeout(30*time.Second))
	testNotImplemetedMethods(t, helper)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cr := containerregistry.NewContainerRegistryClient(helper.authorizer, serverURL, "")
	scopes := containerregistry.AccessTokenScopes{
		{
			Type:    "registry",
			Name:    "catalog",
			Actions: []string{"*"},
		},
	}
	at, claims, err := cr.ExchangeAccessToken(ctx, password, scopes)
	testExpectNoError(t, err)

	if len(at) == 0 {
		t.Fatalf("received empty access token")
	}

	if claims.Issuer != "Azure Container Registry" {
		t.Fatalf("expected the Issuer 'Azure Container Registry' but received: %s", claims.Issuer)
	}

	if claims.Audience != serverURL {
		t.Fatalf("expected the audience %q but received: %s", serverURL, claims.Audience)
	}
}

func testACRHelperGet(t *testing.T, serverURL string, setters ...Option) (*ACRHelper, string) {
	t.Helper()

	helper := NewACRHelper(setters...)
	testNotImplemetedMethods(t, helper)

	username, password, err := helper.Get(serverURL)
	testExpectNoError(t, err)

	if username != "<token>" {
		t.Fatalf("expected username to be '<token>' but received: %s", username)
	}

	if len(password) == 0 {
		t.Fatal("received empty password")
	}

	acrhlp, ok := helper.(*ACRHelper)
	if !ok {
		t.Fatalf("expected helper to be of type '*ACRHelper' but received: %T", helper)
	}

	return acrhlp, password
}

func testNotImplemetedMethods(t *testing.T, helper credentials.Helper) {
	t.Helper()

	err := helper.Add(nil)
	testExpectError(t, err, "method", "not implemented")

	err = helper.Delete("")
	testExpectError(t, err, "method", "not implemented")

	_, err = helper.List()
	testExpectError(t, err, "method", "not implemented")
}

func testExpectError(t *testing.T, err error, containsSubstrings ...string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error but received nil")
	}

	for _, substr := range containsSubstrings {
		if !strings.Contains(err.Error(), substr) {
			t.Fatalf("expected error to contain %q but received: %s", substr, err.Error())
		}
	}
}

func testExpectNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error but received: %v", err)
	}
}
