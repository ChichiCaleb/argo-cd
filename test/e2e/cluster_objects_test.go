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
	"github.com/argoproj/argo-cd/v2/util/argo"
)

func TestClusterRoleBinding(t *testing.T) {
	Given(t).
		Path("cluster-role").
		When().
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced)).
		And(func(app *Application) {
			diffOutput, err := RunCli("app", "diff", app.Name, "--revision=HEAD")
			require.NoError(t, err)
			assert.Empty(t, diffOutput)
		}).
		When().
		SetTrackingMethod(string(argo.TrackingMethodAnnotation)).
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		And(func(app *Application) {
			diffOutput, err := RunCli("app", "diff", app.Name, "--revision=HEAD")
			require.NoError(t, err)
			assert.Empty(t, diffOutput)
		})
}

// ensure that cluster scoped objects, like a cluster role, as a hook, can be successfully deployed
func TestClusterRoleBindingHook(t *testing.T) {
	Given(t).
		Path("cluster-role-hook").
		When().
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced))
}
