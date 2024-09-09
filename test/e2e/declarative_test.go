package e2e

import (
	"testing"

	"github.com/argoproj/gitops-engine/pkg/health"
	. "github.com/argoproj/gitops-engine/pkg/sync/common"

	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	appFixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"
)

func TestDeclarativeHappyApp(t *testing.T) {
	appFixture.Given(t).
		Path("guestbook").
		When().
		Declarative("declarative-apps/app.yaml").
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.HealthIs(health.HealthStatusMissing)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeOutOfSync)).
		When().
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded)).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeSynced))
}

func TestDeclarativeInvalidPath(t *testing.T) {
	appFixture.Given(t).
		Path("garbage").
		When().
		Declarative("declarative-apps/app.yaml").
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeUnknown)).
		Expect(appFixture.Condition(ApplicationConditionComparisonError, "garbage: app path does not exist")).
		When().
		Delete(false).
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.DoesNotExist())
}

func TestDeclarativeInvalidProject(t *testing.T) {
	appFixture.Given(t).
		Path("guestbook").
		Project("garbage").
		When().
		Declarative("declarative-apps/app.yaml").
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.HealthIs(health.HealthStatusUnknown)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeUnknown)).
		Expect(appFixture.Condition(ApplicationConditionInvalidSpecError, "Application referencing project garbage which does not exist"))

	// TODO: you can`t delete application with invalid project due to enforcment that was recently added,
	// in https://github.com/argoproj/argo-cd/security/advisories/GHSA-2gvw-w6fj-7m3c
	// When().
	// Delete(false).
	// Then().
	// Expect(appFixture.Success("")).
	// Expect(DoesNotExist())
}

func TestDeclarativeInvalidRepoURL(t *testing.T) {
	appFixture.Given(t).
		Path("whatever").
		When().
		DeclarativeWithCustomRepo("declarative-apps/app.yaml", "https://github.com").
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.HealthIs(health.HealthStatusHealthy)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeUnknown)).
		Expect(appFixture.Condition(ApplicationConditionComparisonError, "repository not found")).
		When().
		Delete(false).
		Then().
		Expect(appFixture.Success("")).
		Expect(appFixture.DoesNotExist())
}
