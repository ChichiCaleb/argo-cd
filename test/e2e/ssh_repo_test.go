package e2e

import (
	"testing"

	. "github.com/argoproj/gitops-engine/pkg/sync/common"

	"github.com/argoproj/argo-cd/v2/test/e2e/fixture"
	appFixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"
)

func TestCanAccessInsecureSSHRepo(t *testing.T) {
	appFixture.Given(t).
		SSHInsecureRepoURLAdded(true).
		RepoURLType(fixture.RepoURLTypeSSH).
		Path("config-map").
		When().
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded))
}

func TestCanAccessSSHRepo(t *testing.T) {
	appFixture.Given(t).
		CustomSSHKnownHostsAdded().
		SSHRepoURLAdded(true).
		RepoURLType(fixture.RepoURLTypeSSH).
		Path("config-map").
		When().
		CreateApp().
		Sync().
		Then().
		Expect(appFixture.OperationPhaseIs(OperationSucceeded))
}
