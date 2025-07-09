package app

const (
	DeploymentEnvironment_ENVKEY = "ABODEMINE_DEPLOYMENT_ENVIRONMENT"

	// The default environment, just in case.
	DeploymentEnvironment_PRODUCTION = 0

	// A remote environment that is used by clients to test our services
	// without affecting production data.
	// It MUST be updated together with the production environment,
	// and behave as close as possible to it.
	DeploymentEnvironment_SANDBOX = 5

	// A remote environment that can be used by clients to check
	// upcoming features, changes, etc.
	// It SHOULD be updated often with the latest changes before
	// they are pushed to production.
	DeploymentEnvironment_STAGING = 10

	// A remote environment for Dev/QA to test newly developed features.
	// Can be used for stress testing the infrastructure.
	DeploymentEnvironment_TESTING = 20

	// A remote environment that is used for automated tests,
	// usually when changes are pushed to code repositories.
	DeploymentEnvironment_CI = 30

	// A local environment that is used for automated tests,
	// usually to ensure tests pass before pushing changes to code repositories.
	DeploymentEnvironment_LOCAL_CI = 35

	// A local environment that is used by Dev/QA to develop applications.
	DeploymentEnvironment_LOCAL = 40
)
