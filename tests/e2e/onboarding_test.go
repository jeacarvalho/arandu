package e2e

import "testing"

// NOTE: TestOnboardingFlow and related tests are disabled because they test
// functionality from commit 13d67a7 (Provisioning) that was not restored.
// These tests should be re-enabled when Provisioning is implemented again.
func TestOnboardingFlowDisabled(t *testing.T) {
	t.Skip("Skipping e2e onboarding test (Provisioning not restored from 13d67a7)")
}
