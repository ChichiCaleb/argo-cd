package e2e

import (
	"testing"

	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	. "github.com/argoproj/argo-cd/v2/test/e2e/fixture"
	appFixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"

	"github.com/argoproj/gitops-engine/pkg/health"
	. "github.com/argoproj/gitops-engine/pkg/sync/common"

	v1 "k8s.io/api/core/v1"
)

func TestFixingDegradedApp(t *testing.T) {
	appFixture.Given(t).
		Path("sync-waves").
		When().
		IgnoreErrors().
		CreateApp().
		And(func() {
			SetResourceOverrides(map[string]ResourceOverride{
				"ConfigMap": {
					HealthLua: `return { status = obj.metadata.annotations and obj.metadata.annotations['health'] or 'Degraded' }`,
				},
			})
		}).
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationFailed)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeOutOfSync)).
		Expect(appFixture.HealthIs(health.HealthStatusDegraded)).
		Expect(appFixture.ResourceResultNumbering(1)).
		Expect(appFixture.ResourceSyncStatusIs("ConfigMap", "cm-1", SyncStatusCodeSynced)).
		Expect(appFixture.ResourceHealthIs("ConfigMap", "cm-1", health.HealthStatusDegraded)).
		Expect(appFixture.ResourceSyncStatusIs("ConfigMap", "cm-2", SyncStatusCodeOutOfSync)).
		Expect(appFixture.ResourceHealthIs("ConfigMap", "cm-2", health.HealthStatusMissing)).
		When().
		PatchFile("cm-1.yaml", `[{"op": "replace", "path": "/metadata/annotations/health", "value": "Healthy"}]`).
		PatchFile("cm-2.yaml", `[{"op": "replace", "path": "/metadata/annotations/health", "value": "Healthy"}]`).
		// need to force a refresh here
		Refresh(RefreshTypeNormal).
		Then().
		Expect(appFixture.ResourceSyncStatusIs("ConfigMap", "cm-1", SyncStatusCodeOutOfSync)).
		When().
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		Expect(appFixture.ResourceResultNumbering(2)).
		Expect(appFixture.ResourceSyncStatusIs("ConfigMap", "cm-1", SyncStatusCodeSynced)).
		Expect(appFixture.ResourceHealthIs("ConfigMap", "cm-1", health.HealthStatusHealthy)).
		Expect(appFixture.ResourceSyncStatusIs("ConfigMap", "cm-2", SyncStatusCodeSynced)).
		Expect(appFixture.ResourceHealthIs("ConfigMap", "cm-2", health.HealthStatusHealthy))
}

func TestOneProgressingDeploymentIsSucceededAndSynced(t *testing.T) {
	appFixture.Given(t).
		Path("one-deployment").
		When().
		// make this deployment get stuck in progressing due to "invalidimagename"
		PatchFile("deployment.yaml", `[
    {
        "op": "replace",
        "path": "/spec/template/spec/containers/0/image",
        "value": "alpine:ops!"
    }
]`).
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.HealthIs(health.HealthStatusProgressing)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced)).
		Expect(appFixture.ResourceResultNumbering(1))
}

func TestDegradedDeploymentIsSucceededAndSynced(t *testing.T) {
	appFixture.Given(t).
		Path("one-deployment").
		When().
		// make this deployment get stuck in progressing due to "invalidimagename"
		PatchFile("deployment.yaml", `[
    {
        "op": "replace",
        "path": "/spec/progressDeadlineSeconds",
        "value": 1
    },
    {
        "op": "replace",
        "path": "/spec/template/spec/containers/0/image",
        "value": "alpine:ops!"
    }
]`).
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.HealthIs(health.HealthStatusDegraded)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced)).
		Expect(appFixture.ResourceResultNumbering(1))
}

// resources should be pruned in reverse of creation order(syncwaves order)
func TestSyncPruneOrderWithSyncWaves(t *testing.T) {
	ctx := appFixture.Given(t).Timeout(60)

	// remove finalizer to ensure proper cleanup if test fails at early stage
	defer func() {
		_, _ = RunCli("app", "patch-resource", ctx.AppQualifiedName(),
			"--kind", "Pod",
			"--resource-name", "pod-with-finalizers",
			"--patch", `[{"op": "remove", "path": "/metadata/finalizers"}]`,
			"--patch-type", "application/json-patch+json", "--all",
		)
	}()

	ctx.Path("syncwaves-prune-order").
		When().
		CreateApp().
		// creation order: sa & role -> rolebinding -> pod
		Sync().
		Wait().
		Then().
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		When().
		// delete files to remove resources
		DeleteFile("pod.yaml").
		DeleteFile("rbac.yaml").
		Refresh(RefreshTypeHard).
		IgnoreErrors().
		Then().
		Expect(appFixture.SyncStatusIs(SyncStatusCodeOutOfSync)).
		When().
		// prune order: pod -> rolebinding -> sa & role
		Sync("--prune").
		Wait().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		Expect(appFixture.NotPod(func(p v1.Pod) bool { return p.Name == "pod-with-finalizers" })).
		Expect(appFixture.ResourceResultNumbering(4))
}
