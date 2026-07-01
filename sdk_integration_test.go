package edgecron_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	edgecron "github.com/edgecron/edgecron-go"
)

func newIntegrationClient(t *testing.T) *edgecron.Client {
	t.Helper()
	keyID := os.Getenv("EDGECRON_KEY_ID")
	secret := os.Getenv("EDGECRON_SECRET")
	baseURL := os.Getenv("EDGECRON_BASE_URL")
	if keyID == "" || secret == "" {
		t.Skip("EDGECRON_KEY_ID and EDGECRON_SECRET not set")
	}
	if baseURL == "" {
		baseURL = "http://localhost:8888"
	}
	return edgecron.New(keyID, secret, edgecron.WithBaseURL(baseURL))
}

func TestSDKIntegration(t *testing.T) {
	client := newIntegrationClient(t)
	ctx := context.Background()

	// ─── Schedules ───────────────────────────────────────────────────────────

	t.Run("Schedules", func(t *testing.T) {
		name := fmt.Sprintf("test-sched-%d", time.Now().Unix())
		s, err := client.Schedules.Create(ctx, &edgecron.CreateScheduleRequest{
			Name:     name,
			CronExpr: "0 */5 * * *",
			Timezone: "UTC",
			Payload:  `{"key":"val"}`,
		})
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		if s.Name != name {
			t.Fatalf("expected name %q, got %q", name, s.Name)
		}

		got, err := client.Schedules.Get(ctx, s.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.ID != s.ID {
			t.Fatalf("get: expected id %d, got %d", s.ID, got.ID)
		}

		list, err := client.Schedules.List(ctx, 1, 10, "")
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if list.Total < 1 {
			t.Fatalf("list: expected at least 1, got %d", list.Total)
		}

		updatedName := name + "-updated"
		updated, err := client.Schedules.Update(ctx, s.ID, &edgecron.UpdateScheduleRequest{
			Name: &updatedName,
		})
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		if updated.Name != updatedName {
			t.Fatalf("update: expected name %q, got %q", updatedName, updated.Name)
		}

		if err := client.Schedules.Pause(ctx, s.ID); err != nil {
			t.Fatalf("pause: %v", err)
		}

		if err := client.Schedules.Resume(ctx, s.ID); err != nil {
			t.Fatalf("resume: %v", err)
		}

		if err := client.Schedules.Delete(ctx, s.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}
	})

	// ─── Endpoints ──────────────────────────────────────────────────────────

	t.Run("Endpoints", func(t *testing.T) {
		name := fmt.Sprintf("test-ep-%d", time.Now().Unix())
		ep, err := client.Endpoints.Create(ctx, &edgecron.CreateEndpointRequest{
			Name:      name,
			URL:       "https://httpbin.org/post",
			Method:    "POST",
			TimeoutMs: 5000,
		})
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		if ep.Name != name {
			t.Fatalf("expected name %q, got %q", name, ep.Name)
		}

		got, err := client.Endpoints.Get(ctx, ep.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.ID != ep.ID {
			t.Fatalf("get: expected id %d, got %d", ep.ID, got.ID)
		}

		list, err := client.Endpoints.List(ctx, 1, 10, "")
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if list.Total < 1 {
			t.Fatalf("list: expected at least 1, got %d", list.Total)
		}

		updatedName := name + "-updated"
		updated, err := client.Endpoints.Update(ctx, ep.ID, &edgecron.UpdateEndpointRequest{
			Name: &updatedName,
		})
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		if updated.Name != updatedName {
			t.Fatalf("update: expected name %q, got %q", updatedName, updated.Name)
		}

		if err := client.Endpoints.Disable(ctx, ep.ID); err != nil {
			t.Fatalf("disable: %v", err)
		}

		if err := client.Endpoints.Enable(ctx, ep.ID); err != nil {
			t.Fatalf("enable: %v", err)
		}

		if err := client.Endpoints.Delete(ctx, ep.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}
	})

	// ─── Tasks ──────────────────────────────────────────────────────────────

	t.Run("Tasks", func(t *testing.T) {
		endpointID := discoverEndpointID(t, client)
		task, err := client.Tasks.Create(ctx, &edgecron.CreateTaskRequest{
			EndpointID: endpointID,
			Payload:    `{"data":"hello"}`,
		})
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		if task.ID == 0 {
			t.Fatal("create: expected non-zero id")
		}

		got, err := client.Tasks.Get(ctx, task.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.ID != task.ID {
			t.Fatalf("get: expected id %d, got %d", task.ID, got.ID)
		}

		list, err := client.Tasks.List(ctx, 1, 10, "", 0, 0)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if list.Total < 1 {
			t.Fatalf("list: expected at least 1, got %d", list.Total)
		}

		if err := client.Tasks.Cancel(ctx, task.ID); err != nil {
			t.Fatalf("cancel: %v", err)
		}
	})

	// ─── Events ─────────────────────────────────────────────────────────────

	t.Run("Events", func(t *testing.T) {
		eventKey := fmt.Sprintf("test_%d", time.Now().Unix())
		pub, err := client.Events.Publish(ctx, &edgecron.PublishEventRequest{
			EventName: "test",
			EventKey:  eventKey,
			Payload:   `{"order_id":"12345"}`,
		})
		if err != nil {
			t.Fatalf("publish: %v", err)
		}
		if pub.ID == 0 {
			t.Fatal("publish: expected non-zero id")
		}

		got, err := client.Events.Get(ctx, pub.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.ID != pub.ID {
			t.Fatalf("get: expected id %d, got %d", pub.ID, got.ID)
		}

		list, err := client.Events.List(ctx, 1, 10, "", "")
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if list.Total < 1 {
			t.Fatalf("list: expected at least 1, got %d", list.Total)
		}

		if err := client.Events.Disable(ctx, pub.ID); err != nil {
			t.Fatalf("disable: %v", err)
		}

		if err := client.Events.Enable(ctx, pub.ID); err != nil {
			t.Fatalf("enable: %v", err)
		}

		if err := client.Events.Delete(ctx, pub.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}
	})

	// ─── Deliveries ─────────────────────────────────────────────────────────

	t.Run("Deliveries", func(t *testing.T) {
		list, err := client.Deliveries.List(ctx, 1, 10, "", 0, 0)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if list.Total < 0 {
			t.Fatal("list: expected non-negative total")
		}

		// Retry first failed delivery if any
		failed, err := client.Deliveries.List(ctx, 1, 5, "failed", 0, 0)
		if err != nil {
			t.Fatalf("list failed: %v", err)
		}
		for _, d := range failed.List {
			result, err := client.Deliveries.Retry(ctx, d.ID)
			if err != nil {
				t.Logf("retry delivery %d skipped: %v", d.ID, err)
				continue
			}
			if result.DeliveryID != d.ID {
				t.Fatalf("retry: expected delivery_id %d, got %d", d.ID, result.DeliveryID)
			}
			break // only retry one
		}
	})

	// ─── Retries ────────────────────────────────────────────────────────────

	t.Run("Retries", func(t *testing.T) {
		policyName := fmt.Sprintf("test-policy-%d", time.Now().Unix())
		p, err := client.Retries.CreatePolicy(ctx, &edgecron.CreateRetryPolicyRequest{
			Name:            policyName,
			MaxAttempts:     5,
			BackoffType:     "exponential",
			InitialDelaySec: 10,
			MaxDelaySec:     3600,
		})
		if err != nil {
			t.Fatalf("createPolicy: %v", err)
		}
		if p.Name != policyName {
			t.Fatalf("expected name %q, got %q", policyName, p.Name)
		}

		got, err := client.Retries.GetPolicy(ctx, p.ID)
		if err != nil {
			t.Fatalf("getPolicy: %v", err)
		}
		if got.ID != p.ID {
			t.Fatalf("getPolicy: expected id %d, got %d", p.ID, got.ID)
		}

		policies, err := client.Retries.ListPolicies(ctx)
		if err != nil {
			t.Fatalf("listPolicies: %v", err)
		}
		if policies.Total < 1 {
			t.Fatalf("listPolicies: expected at least 1, got %d", policies.Total)
		}

		newMax := int32(10)
		updated, err := client.Retries.UpdatePolicy(ctx, p.ID, &edgecron.UpdateRetryPolicyRequest{
			MaxAttempts: &newMax,
		})
		if err != nil {
			t.Fatalf("updatePolicy: %v", err)
		}
		if updated.MaxAttempts != 10 {
			t.Fatalf("updatePolicy: expected max_attempts 10, got %d", updated.MaxAttempts)
		}

		jobs, err := client.Retries.ListJobs(ctx, 1, 10, "", 0)
		if err != nil {
			t.Fatalf("listJobs: %v", err)
		}
		if jobs.Total < 0 {
			t.Fatal("listJobs: expected non-negative total")
		}

		if err := client.Retries.DeletePolicy(ctx, p.ID); err != nil {
			t.Fatalf("deletePolicy: %v", err)
		}
	})

	// ─── Subscription ───────────────────────────────────────────────────────

	t.Run("Subscription", func(t *testing.T) {
		quota, err := client.Subscription.Quota(ctx)
		if err != nil {
			t.Fatalf("quota: %v", err)
		}
		if quota.PlanCode == "" {
			t.Fatal("quota: expected non-empty plan_code")
		}

		usage, err := client.Subscription.Usage(ctx, "")
		if err != nil {
			t.Fatalf("usage: %v", err)
		}
		if usage.Period == "" {
			t.Fatal("usage: expected non-empty period")
		}

		limits, err := client.Subscription.ResourceLimits(ctx)
		if err != nil {
			t.Fatalf("resourceLimits: %v", err)
		}
		if limits.LogRetentionDays == 0 && limits.MaxCronJobs == 0 {
			t.Fatal("resourceLimits: expected some limits")
		}
	})
}

// discoverEndpointID returns the first endpoint on the account.
func discoverEndpointID(t *testing.T, client *edgecron.Client) int64 {
	t.Helper()
	ctx := context.Background()
	list, err := client.Endpoints.List(ctx, 1, 1, "")
	if err != nil {
		t.Fatalf("discover endpoint: %v", err)
	}
	if len(list.List) == 0 {
		t.Fatal("no endpoints found — create one first")
	}
	return list.List[0].ID
}
