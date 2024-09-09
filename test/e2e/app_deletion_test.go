package e2e

import (
	"testing"

	. "github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/stretchr/testify/assert"

	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	. "github.com/argoproj/argo-cd/v2/test/e2e/fixture"
	appFixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"
)

// when a app gets stuck in sync, and we try to delete it, it won't delete, instead we must then terminate it
// and deletion will then just happen
func TestDeletingAppStuckInSync(t *testing.T) {
	appFixture.Given(t).
		And(func() {
			SetResourceOverrides(map[string]ResourceOverride{
				"ConfigMap": {
					HealthLua: `return { status = obj.annotations and obj.annotations['health'] or 'Progressing' }`,
				},
			})
		}).
		Async(true).
		Path("hook-custom-health").
		When().
		CreateApp().
		Sync().
		Then().
		// stuck in running state
		Expect(appFixture.OperationPhaseIs(OperationRunning)).
		Expect(appFixture.SyncStatusIs(SyncStatusCodeOutOfSync)).
		When().
		Delete(true).
		Then().
		// delete is ignored, still stuck in running state
		Expect(appFixture.OperationPhaseIs(OperationRunning)).
		When().
		TerminateOp().
		Then().
		// delete is successful
		Expect(appFixture.DoesNotExist())
}

func TestDeletingAppByLabel(t *testing.T) {
	appFixture.Given(t).
		Path(guestbookPath).
		When().
		CreateApp("--label=foo=bar").
		Sync().
		Then().
		Expect(appFixture.SyncStatusIs(SyncStatusCode(SyncStatusCodeSynced))).
		When().
		IgnoreErrors().
		DeleteBySelector("foo=baz").
		Then().
		// delete is unsuccessful since no selector match
		AndCLIOutput(
			func(output string, err error) {
				assert.Contains(t, err.Error(), "no apps match selector foo=baz")
			},
		).
		When().
		DeleteBySelector("foo=bar").
		Then().
		// delete is successful
		Expect(appFixture.DoesNotExist())
}

func TestDeletingAppByLabelWait(t *testing.T) {
	appFixture.Given(t).
		Path(guestbookPath).
		When().
		CreateApp("--label=foo=bar").
		Sync().
		Then().
		Expect(appFixture.SyncStatusIs(SyncStatusCode(SyncStatusCodeSynced))).
		When().
		DeleteBySelectorWithWait("foo=bar").
		Then().
		// delete is successful
		Expect(appFixture.DoesNotExistNow())
}
