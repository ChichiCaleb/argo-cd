package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/argoproj/gitops-engine/pkg/health"
	. "github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/stretchr/testify/require"

	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/test/e2e/fixture"
	. "github.com/argoproj/argo-cd/v2/test/e2e/fixture"
	appFixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"
	. "github.com/argoproj/argo-cd/v2/util/errors"
	"github.com/argoproj/argo-cd/v2/util/rand"
)

// when you selectively sync, only selected resources should be synced, but the app will be out of sync
func TestSelectiveSync(t *testing.T) {
	appFixture.Given(t).
		Path("guestbook").
		SelectedResource(":Service:guestbook-ui").
		When().
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeOutOfSync)).
		Expect(appFixture.ResourceHealthIs("Service", "guestbook-ui", health.HealthStatusHealthy)).
		Expect(appFixture.ResourceHealthIs("Deployment", "guestbook-ui", health.HealthStatusMissing))
}

// when running selective sync, hooks do not run
// hooks don't run even if all resources are selected
func TestSelectiveSyncDoesNotRunHooks(t *testing.T) {
	appFixture.Given(t).
		Path("hook").
		SelectedResource(":Pod:pod").
		When().
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		Expect(appFixture.ResourceHealthIs("Pod", "pod", health.HealthStatusHealthy)).
		Expect(appFixture.ResourceResultNumbering(1))
}

func TestSelectiveSyncWithoutNamespace(t *testing.T) {
	selectedResourceNamespace := getNewNamespace(t)
	defer func() {
		if !t.Skipped() {
			FailOnErr(Run("", "kubectl", "delete", "namespace", selectedResourceNamespace))
		}
	}()
	appFixture.Given(t).
		Prune(true).
		Path("guestbook-with-namespace").
		And(func() {
			FailOnErr(Run("", "kubectl", "create", "namespace", selectedResourceNamespace))
		}).
		SelectedResource("apps:Deployment:guestbook-ui").
		When().
		PatchFile("guestbook-ui-deployment-ns.yaml", fmt.Sprintf(`[{"op": "replace", "path": "/metadata/namespace", "value": "%s"}]`, selectedResourceNamespace)).
		PatchFile("guestbook-ui-svc-ns.yaml", fmt.Sprintf(`[{"op": "replace", "path": "/metadata/namespace", "value": "%s"}]`, selectedResourceNamespace)).
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeOutOfSync)).
		Expect(appFixture.ResourceHealthWithNamespaceIs("Deployment", "guestbook-ui", selectedResourceNamespace, health.HealthStatusHealthy)).
		Expect(appFixture.ResourceHealthWithNamespaceIs("Deployment", "guestbook-ui", fixture.DeploymentNamespace(), health.HealthStatusHealthy)).
		Expect(appFixture.ResourceSyncStatusWithNamespaceIs("Deployment", "guestbook-ui", selectedResourceNamespace, SyncStatusCodeSynced)).
		Expect(appFixture.ResourceSyncStatusWithNamespaceIs("Deployment", "guestbook-ui", fixture.DeploymentNamespace(), SyncStatusCodeSynced))
}

// In selectedResource to sync, namespace is provided
func TestSelectiveSyncWithNamespace(t *testing.T) {
	selectedResourceNamespace := getNewNamespace(t)
	defer func() {
		if !t.Skipped() {
			FailOnErr(Run("", "kubectl", "delete", "namespace", selectedResourceNamespace))
		}
	}()
	appFixture.Given(t).
		Prune(true).
		Path("guestbook-with-namespace").
		And(func() {
			FailOnErr(Run("", "kubectl", "create", "namespace", selectedResourceNamespace))
		}).
		SelectedResource(fmt.Sprintf("apps:Deployment:%s/guestbook-ui", selectedResourceNamespace)).
		When().
		PatchFile("guestbook-ui-deployment-ns.yaml", fmt.Sprintf(`[{"op": "replace", "path": "/metadata/namespace", "value": "%s"}]`, selectedResourceNamespace)).
		PatchFile("guestbook-ui-svc-ns.yaml", fmt.Sprintf(`[{"op": "replace", "path": "/metadata/namespace", "value": "%s"}]`, selectedResourceNamespace)).
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeOutOfSync)).
		Expect(appFixture.ResourceHealthWithNamespaceIs("Deployment", "guestbook-ui", selectedResourceNamespace, health.HealthStatusHealthy)).
		Expect(appFixture.ResourceHealthWithNamespaceIs("Deployment", "guestbook-ui", fixture.DeploymentNamespace(), health.HealthStatusMissing)).
		Expect(appFixture.ResourceSyncStatusWithNamespaceIs("Deployment", "guestbook-ui", selectedResourceNamespace, SyncStatusCodeSynced)).
		Expect(appFixture.ResourceSyncStatusWithNamespaceIs("Deployment", "guestbook-ui", fixture.DeploymentNamespace(), SyncStatusCodeOutOfSync))
}

func getNewNamespace(t *testing.T) string {
	randStr, err := rand.String(5)
	require.NoError(t, err)
	postFix := "-" + strings.ToLower(randStr)
	name := fixture.DnsFriendly(t.Name(), "")
	return fixture.DnsFriendly(fmt.Sprintf("argocd-e2e-%s", name), postFix)
}
