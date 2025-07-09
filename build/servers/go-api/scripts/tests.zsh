source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Function to check if a file exists and is readable
file_exists() {
  [ -f "$1" ] && [ -r "$1" ]
}

# Function to validate JSON
validate_json() {
  echo "$1" | jq '.' >/dev/null 2>&1
}

# Extract Postman test parameters from the build params JSON file
if [ ! -f "${ABODEMINE_BUILD_PARAMS}" ]; then
  echo "::error::Build params file not found: ${ABODEMINE_BUILD_PARAMS}"
  exit 1
fi

# First try to read from local environment variable, if set
if [ -n "${POSTMAN_ENVIRONMENT_ID}" ]; then
  ENVIRONMENT_ID="${POSTMAN_ENVIRONMENT_ID}"
else
  # Fall back to reading from build params
  ENVIRONMENT_ID=$(jq -r '."projects/'${ABODEMINE_PROJECT_SLUG}'".tests.postman.environment_id' ${ABODEMINE_BUILD_PARAMS})
fi
API_KEY=$(jq -r '."projects/'${ABODEMINE_PROJECT_SLUG}'".tests.postman.api_key' ${ABODEMINE_BUILD_PARAMS})
PROJECT_API_BEARER_TOKEN=$(jq -r '."projects/'${ABODEMINE_PROJECT_SLUG}'".tests.postman.project_api_bearer_token' ${ABODEMINE_BUILD_PARAMS})

# Validate extracted values
if [ -z "${ENVIRONMENT_ID}" ] || [ "${ENVIRONMENT_ID}" = "null" ]; then
  echo "::error::Invalid or missing Postman Environment ID"
  exit 1
fi

if [ -z "${PROJECT_API_BEARER_TOKEN}" ] || [ "${PROJECT_API_BEARER_TOKEN}" = "null" ]; then
  echo "::error::Invalid or missing Postman Project API Bearer Token"
  exit 1
fi

if [ -z "${API_KEY}" ] || [ "${API_KEY}" = "null" ]; then
  echo "::error::Invalid or missing Postman API Key"
  exit 1
fi

# Print extracted values for debugging
echo "::group::Postman Configuration"
echo "Postman Environment ID: ${ENVIRONMENT_ID}"
echo "::endgroup::"

# Download the environment file
echo "::group::Downloading Postman Environment"
ENVIRONMENT_FILE="${ABODEMINE_BUILD_TMP}/postman_environment.json"
curl -s -H "X-API-Key: ${API_KEY}" \
  "https://api.getpostman.com/environments/${ENVIRONMENT_ID}" | jq '.environment' > ${ENVIRONMENT_FILE}
echo "::endgroup::"

echo "::group::Running Postman Tests"
# Check if newman is installed, if not install it
if ! command -v newman &> /dev/null; then   
  echo "Newman not found, installing..."
  npm install -g newman newman-reporter-htmlextra
fi

# Run Newman with enhanced configuration
# Using set +e to prevent script from exiting if newman fails
set +e
newman run "${ABODEMINE_WORKSPACE}/code/json/tests/${ABODEMINE_PROJECT_SLUG}.postman_collection.json" \
  -e "${ENVIRONMENT_FILE}" \
  --insecure \
  --env-var "projectApiBearerToken=${PROJECT_API_BEARER_TOKEN}" \
  -r cli,junit,htmlextra \
  --reporter-junit-export "${ABODEMINE_WORKSPACE}/build/servers/go-api/newman_junit_output.xml" \
  --reporter-htmlextra-export "${ABODEMINE_WORKSPACE}/build/servers/go-api/newman_htmlextra_output.html"
NEWMAN_EXIT_CODE=$?
# Restore error handling after newman command
set -e

# Print test summary
echo "::endgroup::"
echo "::group::Test Summary"
if [ -f "${ABODEMINE_WORKSPACE}/build/servers/go-api/newman_junit_output.xml" ]; then
  echo "JUnit report generated successfully"
  echo "Report location: ${ABODEMINE_WORKSPACE}/build/servers/go-api/newman_junit_output.xml"
else
  echo "::error::Failed to generate JUnit report"
fi

if [ ${NEWMAN_EXIT_CODE} -eq 0 ]; then
  echo "::notice::All tests passed successfully"
else
  echo "::error::Some tests failed"
fi
echo "::endgroup::"

exit ${NEWMAN_EXIT_CODE}