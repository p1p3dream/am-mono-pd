package lambda

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/lib/must"
	"abodemine/projects/datapipe/conf"
)

/*
GO_TEST_COUNT=1 \
GO_TEST_PARAMS="-v" \
make -C ${ABODEMINE_WORKSPACE}/code/go/abodemine \
test/abodemine/projects/datapipe/domains/lambda \
RUN="TestDomain_HandleLambdaEventInput"
*/
func TestDomain_HandleLambdaEventInput(t *testing.T) {
	if os.Getenv("RUN") != t.Name() {
		t.Skip("This test MUST be selected manually.")
	}

	testCases := []*struct {
		name string
		in   *HandleTaskLauncherLambdaEventInput
		out  *HandleTaskLauncherLambdaEventOutput
		err  *errors.Object
	}{
		{
			name: "nil-input",
			err: &errors.Object{
				Id:   "687f5932-835d-4ee9-8a9f-f52d58830639",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "ok",
			in:   &HandleTaskLauncherLambdaEventInput{},
		},
	}

	ctx := context.Background()
	config := conf.MustResolveAndLoadOnce(ctx)

	requestDomain := arc.NewDomain(&arc.NewDomainInput{})
	workerDomain := NewDomain(&NewDomainInput{
		Config: config,
	})

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := requestDomain.CreateRequest(&arc.CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatal(err)
			}

			_, err = workerDomain.HandleTaskLauncherLambdaEvent(r, tc.in)
			if err != nil {
				if tc.err == nil {
					st.Fatalf("failed to HandleLambdaEvent: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")

				return
			}
		})
	}
}
