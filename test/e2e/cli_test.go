package e2e

import (
	"testing"

	"github.com/argoproj/gitops-engine/pkg/health"
	. "github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	. "github.com/argoproj/argo-cd/v2/test/e2e/fixture"
	appFixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"
)

func TestCliAppCommand(t *testing.T) {
	Given(t).
		Path("hook").
		When().
		CreateApp().
		And(func() {
			output, err := RunCli("app", "sync", Name(), "--timeout", "90")
			require.NoError(t, err)
			vars := map[string]interface{}{"Name": Name(), "Namespace": DeploymentNamespace()}
			assert.Contains(t, NormalizeOutput(output), Tmpl(`Pod {{.Namespace}} pod Synced Progressing pod/pod created`, vars))
			assert.Contains(t, NormalizeOutput(output), Tmpl(`Pod {{.Namespace}} hook Succeeded Sync pod/hook created`, vars))
		}).
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		And(func(_ *Application) {
			output, err := RunCli("app", "list")
			require.NoError(t, err)
			expected := Tmpl(
				`{{.Name}} https://kubernetes.default.svc {{.Namespace}} default Synced Healthy Manual <none>`,
				map[string]interface{}{"Name": Name(), "Namespace": DeploymentNamespace()})
			assert.Contains(t, NormalizeOutput(output), expected)
		})
}
