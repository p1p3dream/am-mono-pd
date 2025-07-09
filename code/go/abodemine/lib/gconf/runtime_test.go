package gconf

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"abodemine/lib/app"
	"abodemine/lib/unique"
)

func TestResolveDeploymentEnvironment(t *testing.T) {
	// Since this is a critical part of the system,
	// it's better to test it thoroughly.

	testCases := []*struct {
		name string
		in   string
		want int
	}{
		{
			name: "empty-env-value",
			want: app.DeploymentEnvironment_PRODUCTION,
		},
		{
			name: "testing-a",
			in:   "testing-a",
			want: app.DeploymentEnvironment_TESTING,
		},
		{
			name: "local_ci-b",
			in:   "local_ci-b",
			want: app.DeploymentEnvironment_LOCAL_CI,
		},
		{
			name: "unexistent-c",
			in:   "unexistent-c",
			want: app.DeploymentEnvironment_PRODUCTION,
		},
		{
			name: "random",
			in:   uuid.NewString(),
			want: app.DeploymentEnvironment_PRODUCTION,
		},
		{
			name: "production",
			in:   "production",
			want: app.DeploymentEnvironment_PRODUCTION,
		},
		{
			name: "sandbox",
			in:   "sandbox",
			want: app.DeploymentEnvironment_SANDBOX,
		},
		{
			name: "staging",
			in:   "staging",
			want: app.DeploymentEnvironment_STAGING,
		},
		{
			name: "testing",
			in:   "testing",
			want: app.DeploymentEnvironment_TESTING,
		},
		{
			name: "ci",
			in:   "ci",
			want: app.DeploymentEnvironment_CI,
		},
		{
			name: "local_ci",
			in:   "local_ci",
			want: app.DeploymentEnvironment_LOCAL_CI,
		},
		{
			name: "local",
			in:   "local",
			want: app.DeploymentEnvironment_LOCAL,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			envKey := unique.Upper(16)
			os.Setenv(envKey, tc.in)

			have := ResolveDeploymentEnvironment(envKey)

			assert.Equal(st, tc.want, have)
		})
	}
}
