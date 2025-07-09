package gconf

import (
	"os"
	"strings"

	"abodemine/lib/app"
)

// ResolveDeploymentEnvironment resolves the deployment environment based on the
// value of the environment variable with the given key.
// Defaults to app.DeploymentEnvironment_ENVKEY.
func ResolveDeploymentEnvironment(envKey string) int {
	if envKey == "" {
		envKey = app.DeploymentEnvironment_ENVKEY
	}

	v := strings.ToUpper(strings.TrimSpace(os.Getenv(envKey)))

	// Deployment environments can be suffixed with a dash
	// to indicate they are a sub-environment.
	// This means that for business rules they are the same,
	// but probably have different configurations.
	if idx := strings.Index(v, "-"); idx != -1 {
		v = v[:idx]
	}

	return DeploymentEnvironmentFromString(v)
}

func DeploymentEnvironmentFromString(s string) int {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "SANDBOX":
		return app.DeploymentEnvironment_SANDBOX
	case "STAGING":
		return app.DeploymentEnvironment_STAGING
	case "TESTING":
		return app.DeploymentEnvironment_TESTING
	case "CI":
		return app.DeploymentEnvironment_CI
	case "LOCAL_CI":
		return app.DeploymentEnvironment_LOCAL_CI
	case "LOCAL":
		return app.DeploymentEnvironment_LOCAL
	}

	return app.DeploymentEnvironment_PRODUCTION
}

func DeploymentEnvironmentToString(env int) string {
	switch env {
	case app.DeploymentEnvironment_SANDBOX:
		return "SANDBOX"
	case app.DeploymentEnvironment_STAGING:
		return "STAGING"
	case app.DeploymentEnvironment_TESTING:
		return "TESTING"
	case app.DeploymentEnvironment_CI:
		return "CI"
	case app.DeploymentEnvironment_LOCAL_CI:
		return "LOCAL_CI"
	case app.DeploymentEnvironment_LOCAL:
		return "LOCAL"
	}

	return "PRODUCTION"
}
