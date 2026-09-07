package scripts

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

func runRoute(t *testing.T, eventName, ref, refType, eventAction, inputsMode, eventSchedule string) map[string]string {
	cmd := exec.Command("bash", "./ci_route.sh", eventName, ref, refType, eventAction, inputsMode, eventSchedule)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("ci_route.sh failed: %v", err)
	}

	result := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			result[parts[0]] = parts[1]
		}
	}
	// On error, exec.Run returns an error which is caught.
	// We'll capture exit code separately if needed, but for our tests we check stderr or error returned
	return result
}

func runRouteWithError(t *testing.T, eventName, ref, refType, eventAction, inputsMode, eventSchedule string) error {
	cmd := exec.Command("bash", "./ci_route.sh", eventName, ref, refType, eventAction, inputsMode, eventSchedule)
	return cmd.Run()
}

func runSnapshot(t *testing.T, eventName, inputsSnapshotMode, releaseTag, ref string) string {
	cmd := exec.Command("bash", "./ci_snapshot.sh", eventName, inputsSnapshotMode, releaseTag, ref)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("ci_snapshot.sh failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "is_snapshot=") {
			return strings.TrimPrefix(line, "is_snapshot=")
		}
	}
	return ""
}

func runDispatchMode(t *testing.T, mode string) string {
	cmd := exec.Command("bash", "./ci_dispatch_mode.sh", mode)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("ci_dispatch_mode.sh failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "snapshot_mode=") {
			return strings.TrimPrefix(line, "snapshot_mode=")
		}
	}
	return ""
}

func runWindowsUpload(t *testing.T, eventName, inputsSnapshotMode, releaseTag string) string {
	cmd := exec.Command("bash", "./ci_windows_upload.sh", eventName, inputsSnapshotMode, releaseTag)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("ci_windows_upload.sh failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "windows_upload=") {
			return strings.TrimPrefix(line, "windows_upload=")
		}
	}
	return ""
}

func TestPolicyNormalBranchPush(t *testing.T) {
	// 1. Normal branch push
	// - validation/build checks as intended;
	// - no release publication.
	res := runRoute(t, "push", "refs/heads/main", "branch", "", "", "")
	if res["run_code_checks"] != "true" || res["run_release"] != "false" {
		t.Errorf("Expected run_code_checks=true, run_release=false, got %v", res)
	}
}

func TestPolicyPullRequest(t *testing.T) {
	// 2. Pull request
	// - normal PR checks;
	// - no release publication.
	res := runRoute(t, "pull_request", "refs/pull/1/merge", "", "opened", "", "")
	if res["run_code_checks"] != "true" || res["run_pr_meta_checks"] != "true" || res["run_release"] != "false" {
		t.Errorf("Expected run_code_checks=true, run_pr_meta_checks=true, run_release=false, got %v", res)
	}
}

func TestPolicyManualReleasePublish(t *testing.T) {
	// 3. Manual `release-major`, `release-minor`, `release-patch`
	// - route as normal published releases;
	modes := []string{"release-major", "release-minor", "release-patch"}
	for _, mode := range modes {
		res := runRoute(t, "workflow_dispatch", "", "", "", mode, "")
		if res["run_code_checks"] != "true" || res["run_release"] != "true" {
			t.Errorf("Mode %s: Expected run_code_checks=true, run_release=true, got %v", mode, res)
		}

		// - normal release dispatch does not use snapshot_mode=true
		dispatchSnap := runDispatchMode(t, mode)
		if dispatchSnap != "false" {
			t.Errorf("Mode %s: Expected dispatch snapshot_mode=false, got %v", mode, dispatchSnap)
		}
	}
}

func TestPolicyManualReleaseSnapshot(t *testing.T) {
	// 4. Manual `release-test`, `release-rc`, `release-alpha`
	// Note: in the real CI, release-test/rc/alpha trigger another publish-tag workflow with snapshot_mode=true
	// So we need to test that publish-tag + snapshot_mode=true sets snapshot mode correctly
	// The routing logic maps release-* to run_release=true. Snapshot intent is preserved via publish-tag hop.
	modes := []string{"release-test", "release-rc", "release-alpha"}
	for _, mode := range modes {
		res := runRoute(t, "workflow_dispatch", "", "", "", mode, "")
		if res["run_code_checks"] != "true" || res["run_release"] != "true" {
			t.Errorf("Mode %s: Expected run_code_checks=true, run_release=true, got %v", mode, res)
		}

		// - test/rc/alpha dispatch uses snapshot_mode=true
		dispatchSnap := runDispatchMode(t, mode)
		if dispatchSnap != "true" {
			t.Errorf("Mode %s: Expected dispatch snapshot_mode=true, got %v", mode, dispatchSnap)
		}
	}
}

func TestPolicyInternalPublishTag(t *testing.T) {
	// 5. Internal `publish-tag`
	// - validate required tag context;
	// - distinguish snapshot vs normal policy correctly.
	res := runRoute(t, "workflow_dispatch", "refs/tags/v1.0.0", "tag", "", "publish-tag", "")
	if res["run_code_checks"] != "true" || res["run_release"] != "true" {
		t.Errorf("Expected run_code_checks=true, run_release=true, got %v", res)
	}

	snapNormal := runSnapshot(t, "workflow_dispatch", "false", "v1.0.0", "refs/tags/v1.0.0")
	if snapNormal != "false" {
		t.Errorf("Expected is_snapshot=false, got %v", snapNormal)
	}

	snapTest := runSnapshot(t, "workflow_dispatch", "true", "v1.0.0", "refs/tags/v1.0.0")
	if snapTest != "true" {
		t.Errorf("Expected is_snapshot=true for snapshot_mode=true, got %v", snapTest)
	}
}

func TestPolicyInternalPublishTagInvalid(t *testing.T) {
	// - validate required tag context (failure mode)
	err := runRouteWithError(t, "workflow_dispatch", "refs/heads/main", "branch", "", "publish-tag", "")
	if err == nil {
		t.Errorf("Expected an error exit code for invalid ref context, but got success")
	}
}

func TestPolicyExternalVTag(t *testing.T) {
	// 6. External `v*` tag:
	// - normal publication policy.
	res := runRoute(t, "push", "refs/tags/v1.0.0", "tag", "", "", "")
	if res["run_code_checks"] != "true" || res["run_release"] != "true" {
		t.Errorf("Expected run_code_checks=true, run_release=true, got %v", res)
	}
	snap := runSnapshot(t, "push", "false", "v1.0.0", "refs/tags/v1.0.0")
	if snap != "false" {
		t.Errorf("Expected is_snapshot=false, got %v", snap)
	}
}

func TestPolicyExternalTestTag(t *testing.T) {
	// 7. External `test-*` tag:
	// - snapshot-only;
	// - must never become a normal GitHub Release.
	res := runRoute(t, "push", "refs/tags/test-1", "tag", "", "", "")
	if res["run_code_checks"] != "true" || res["run_release"] != "true" {
		t.Errorf("Expected run_code_checks=true, run_release=true, got %v", res)
	}
	snap := runSnapshot(t, "push", "false", "test-1", "refs/tags/test-1")
	if snap != "true" {
		t.Errorf("Expected is_snapshot=true, got %v", snap)
	}
}

func TestPolicyFailedChecks(t *testing.T) {
	// 8. Failed lint/test/vet:
	// - release validation fails or does not pass;
	// - `release-context` cannot push/dispatch a release tag.
	// We statically test this invariant by reading ci.yml and ensuring release-validation
	// requires golangci, go-test, and go-vet, and release-context requires release-validation.
	content, err := ioutil.ReadFile("../.github/workflows/ci.yml")
	if err != nil {
		t.Fatalf("Failed to read ci.yml: %v", err)
	}
	output := strings.ReplaceAll(string(content), "\r\n", "\n")

	if !strings.Contains(output, "name: Release Validation Gate\n    needs: [route, discover, prepare-release-tag, golangci, go-test, go-vet]") {
		t.Errorf("Expected Release Validation Gate to depend on lint and tests")
	}

	if !strings.Contains(output, "name: Release Context & Gate\n    needs: [route, prepare-release-tag, release-validation]") {
		t.Errorf("Expected Release Context & Gate to depend on release-validation")
	}
}

func TestPolicySingleOwner(t *testing.T) {
	// 9. Single owner:
	// - exactly one lane/tool owns GitHub Release creation.
	// - Linux GoReleaser is the canonical publisher under the current architecture.
	// - Windows must not independently publish the same release.
	content, err := ioutil.ReadFile("../.github/workflows/ci.yml")
	if err != nil {
		t.Fatalf("Failed to read ci.yml: %v", err)
	}
	output := strings.ReplaceAll(string(content), "\r\n", "\n")

	if !strings.Contains(output, "release --clean -f .goreleaser-linux.yml") {
		t.Errorf("Expected Linux to run release --clean")
	}

	if !strings.Contains(output, "needs: [route, discover, release-context, goreleaser-linux]") {
		t.Errorf("Expected Windows to depend on goreleaser-linux")
	}
}

func TestPolicyWindowsArtifactUpload(t *testing.T) {
	// 10. Windows artifact upload:
	// - enabled only when a canonical published GitHub Release is supposed to exist;
	// - disabled for snapshot/test policy.
	if runWindowsUpload(t, "push", "false", "v1.0.0") != "true" {
		t.Errorf("Expected windows upload true for normal release")
	}
	if runWindowsUpload(t, "push", "false", "test-1") != "false" {
		t.Errorf("Expected windows upload false for test tag")
	}
	if runWindowsUpload(t, "workflow_dispatch", "true", "v1.0.0") != "false" {
		t.Errorf("Expected windows upload false for snapshot mode")
	}
}

func TestPolicyReleasePublished(t *testing.T) {
	// 11. `release: published`:
	// - downstream-only;
	// - cannot route back into tag/release creation.
	res := runRoute(t, "release", "", "", "published", "", "")
	if res["run_post_release"] != "true" || res["run_release"] != "false" {
		t.Errorf("Expected run_post_release=true and run_release=false, got %v", res)
	}
}
