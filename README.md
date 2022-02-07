# Docker Credential Helper for Azure Container Registry

This is a small project that tries to implement the [docker-credential-helpers](https://github.com/docker/docker-credential-helpers) `credentials.Helper` interface without using the Microsoft SDKs but rather the [manicminer/hamilton](https://github.com/manicminer/hamilton) Go SDK for Microsoft Graph.

This project is currently using a [fork of hamilton](https://github.com/simongottschlag/hamilton/tree/docker-credential-helper) and waiting for the following two PRs to merge:

- [Separate autorest module #139](https://github.com/manicminer/hamilton/pull/139)
- [services/containerregistry: Add support for Azure Container Registry token exchange #144](https://github.com/manicminer/hamilton/pull/144)

Alternatives:

- [chrismellard/docker-credential-acr-env](https://github.com/chrismellard/docker-credential-acr-env)
- [Azure/acr-docker-credential-helper](https://github.com/Azure/acr-docker-credential-helper)
