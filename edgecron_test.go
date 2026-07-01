package edgecron_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	edgecron "github.com/edgecron/edgecron-go"
)

func mockServer(t *testing.T, wantPath, wantMethod string, dataJSON string) (*httptest.Server, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path: got %q, want %q", r.URL.Path, wantPath)
		}
		if r.Method != wantMethod {
			t.Errorf("method: got %q, want %q", r.Method, wantMethod)
		}
		if r.Header.Get("X-Key-ID") == "" {
			t.Error("missing X-Key-ID header")
		}
		if r.Header.Get("X-Timestamp") == "" {
			t.Error("missing X-Timestamp header")
		}
		if r.Header.Get("X-Signature") == "" {
			t.Error("missing X-Signature header")
		}
		resp := map[string]interface{}{"code": 0, "message": "success", "request_id": "test-rid"}
		if dataJSON != "" {
			var data interface{}
			json.Unmarshal([]byte(dataJSON), &data)
			resp["data"] = data
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	return srv, srv.Close
}

func newTestClient(t *testing.T, srv *httptest.Server) *edgecron.Client {
	t.Helper()
	return edgecron.New("ak_3f9a2b1c7d4e8f0a", "test-secret", edgecron.WithBaseURL(srv.URL))
}

// --- Schedules ---

func TestSchedules_Create(t *testing.T) {
	srv, close := mockServer(t, "/v1/schedules", http.MethodPost, `{"id":1,"app_id":"app_1","name":"test","cron_expr":"* * * * *","status":"active"}`)
	defer close()
	client := newTestClient(t, srv)
	res, err := client.Schedules.Create(context.Background(), &edgecron.CreateScheduleRequest{Name: "test", CronExpr: "* * * * *"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "test" || res.Status != "active" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestSchedules_Get(t *testing.T) {
	srv, close := mockServer(t, "/v1/schedules/42", http.MethodGet, `{"id":42,"name":"test","cron_expr":"* * * * *","status":"active"}`)
	defer close()
	res, err := newTestClient(t, srv).Schedules.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != 42 {
		t.Fatalf("expected id 42, got %d", res.ID)
	}
}

func TestSchedules_List(t *testing.T) {
	srv, close := mockServer(t, "/v1/schedules", http.MethodGet, `{"total":0,"list":[]}`)
	defer close()
	_, err := newTestClient(t, srv).Schedules.List(context.Background(), 1, 20, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchedules_Pause(t *testing.T) {
	srv, close := mockServer(t, "/v1/schedules/1/pause", http.MethodPost, "")
	defer close()
	err := newTestClient(t, srv).Schedules.Pause(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchedules_Resume(t *testing.T) {
	srv, close := mockServer(t, "/v1/schedules/1/resume", http.MethodPost, "")
	defer close()
	err := newTestClient(t, srv).Schedules.Resume(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchedules_Delete(t *testing.T) {
	srv, close := mockServer(t, "/v1/schedules/1", http.MethodDelete, "")
	defer close()
	err := newTestClient(t, srv).Schedules.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Tasks ---

func TestTasks_Create(t *testing.T) {
	srv, close := mockServer(t, "/v1/tasks", http.MethodPost, `{"id":1,"task_type":"manual","status":"pending"}`)
	defer close()
	res, err := newTestClient(t, srv).Tasks.Create(context.Background(), &edgecron.CreateTaskRequest{EndpointID: 1, Payload: "{}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != 1 || res.Status != "pending" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestTasks_Get(t *testing.T) {
	srv, close := mockServer(t, "/v1/tasks/5", http.MethodGet, `{"id":5,"task_type":"schedule","status":"running"}`)
	defer close()
	res, err := newTestClient(t, srv).Tasks.Get(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != 5 {
		t.Fatalf("expected id 5, got %d", res.ID)
	}
}

func TestTasks_List(t *testing.T) {
	srv, close := mockServer(t, "/v1/tasks", http.MethodGet, `{"total":0,"list":[]}`)
	defer close()
	_, err := newTestClient(t, srv).Tasks.List(context.Background(), 1, 20, "", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTasks_Cancel(t *testing.T) {
	srv, close := mockServer(t, "/v1/tasks/1/cancel", http.MethodPost, "")
	defer close()
	err := newTestClient(t, srv).Tasks.Cancel(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Events ---

func TestEvents_Publish(t *testing.T) {
	srv, close := mockServer(t, "/v1/events", http.MethodPost, `{"id":1,"event_name":"user.created","event_key":"evt_001","status":"pending","fanout_count":2}`)
	defer close()
	res, err := newTestClient(t, srv).Events.Publish(context.Background(), &edgecron.PublishEventRequest{EventName: "user.created", EventKey: "evt_001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.EventName != "user.created" || res.FanoutCount != 2 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestEvents_Get(t *testing.T) {
	srv, close := mockServer(t, "/v1/events/10", http.MethodGet, `{"id":10,"event_name":"test.event","status":"active"}`)
	defer close()
	res, err := newTestClient(t, srv).Events.Get(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != 10 {
		t.Fatalf("expected id 10, got %d", res.ID)
	}
}

func TestEvents_List(t *testing.T) {
	srv, close := mockServer(t, "/v1/events", http.MethodGet, `{"total":0,"list":[]}`)
	defer close()
	_, err := newTestClient(t, srv).Events.List(context.Background(), 1, 20, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEvents_Enable(t *testing.T) {
	srv, close := mockServer(t, "/v1/events/1/enable", http.MethodPost, "")
	defer close()
	err := newTestClient(t, srv).Events.Enable(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEvents_Disable(t *testing.T) {
	srv, close := mockServer(t, "/v1/events/1/disable", http.MethodPost, "")
	defer close()
	err := newTestClient(t, srv).Events.Disable(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEvents_Delete(t *testing.T) {
	srv, close := mockServer(t, "/v1/events/1", http.MethodDelete, "")
	defer close()
	err := newTestClient(t, srv).Events.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Endpoints ---

func TestEndpoints_Create(t *testing.T) {
	srv, close := mockServer(t, "/v1/endpoints", http.MethodPost, `{"id":1,"name":"my-ep","url":"https://example.com/hook","method":"POST","status":"enabled"}`)
	defer close()
	res, err := newTestClient(t, srv).Endpoints.Create(context.Background(), &edgecron.CreateEndpointRequest{Name: "my-ep", URL: "https://example.com/hook"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "my-ep" || res.Status != "enabled" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestEndpoints_Get(t *testing.T) {
	srv, close := mockServer(t, "/v1/endpoints/3", http.MethodGet, `{"id":3,"name":"test-ep","status":"enabled"}`)
	defer close()
	res, err := newTestClient(t, srv).Endpoints.Get(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != 3 {
		t.Fatalf("expected id 3, got %d", res.ID)
	}
}

func TestEndpoints_Update(t *testing.T) {
	srv, close := mockServer(t, "/v1/endpoints/1", http.MethodPatch, `{"id":1,"name":"updated-ep","url":"https://example.com/new","status":"enabled"}`)
	defer close()
	res, err := newTestClient(t, srv).Endpoints.Update(context.Background(), 1, &edgecron.UpdateEndpointRequest{Name: strPtr("updated-ep"), URL: strPtr("https://example.com/new")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "updated-ep" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestEndpoints_List(t *testing.T) {
	srv, close := mockServer(t, "/v1/endpoints", http.MethodGet, `{"total":0,"list":[]}`)
	defer close()
	_, err := newTestClient(t, srv).Endpoints.List(context.Background(), 1, 20, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEndpoints_Delete(t *testing.T) {
	srv, close := mockServer(t, "/v1/endpoints/1", http.MethodDelete, "")
	defer close()
	err := newTestClient(t, srv).Endpoints.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEndpoints_Enable(t *testing.T) {
	srv, close := mockServer(t, "/v1/endpoints/1/enable", http.MethodPost, "")
	defer close()
	err := newTestClient(t, srv).Endpoints.Enable(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEndpoints_Disable(t *testing.T) {
	srv, close := mockServer(t, "/v1/endpoints/1/disable", http.MethodPost, "")
	defer close()
	err := newTestClient(t, srv).Endpoints.Disable(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Deliveries ---

func TestDeliveries_List(t *testing.T) {
	srv, close := mockServer(t, "/v1/deliveries", http.MethodGet, `{"total":0,"list":[]}`)
	defer close()
	_, err := newTestClient(t, srv).Deliveries.List(context.Background(), 1, 20, "", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeliveries_Retry(t *testing.T) {
	srv, close := mockServer(t, "/v1/deliveries/1/retry", http.MethodPost, `{"delivery_id":1,"retry_job_id":100,"status":"retry_scheduled"}`)
	defer close()
	res, err := newTestClient(t, srv).Deliveries.Retry(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.DeliveryID != 1 || res.RetryJobID != 100 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

// --- Retries ---

func TestRetries_CreatePolicy(t *testing.T) {
	srv, close := mockServer(t, "/v1/retries/policies", http.MethodPost, `{"id":1,"name":"default","max_attempts":3,"backoff_type":"exponential","status":"active"}`)
	defer close()
	res, err := newTestClient(t, srv).Retries.CreatePolicy(context.Background(), &edgecron.CreateRetryPolicyRequest{Name: "default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "default" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestRetries_GetPolicy(t *testing.T) {
	srv, close := mockServer(t, "/v1/retries/policies/7", http.MethodGet, `{"id":7,"name":"my-policy","max_attempts":5,"status":"active"}`)
	defer close()
	res, err := newTestClient(t, srv).Retries.GetPolicy(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != 7 {
		t.Fatalf("expected id 7, got %d", res.ID)
	}
}

func TestRetries_UpdatePolicy(t *testing.T) {
	srv, close := mockServer(t, "/v1/retries/policies/1", http.MethodPatch, `{"id":1,"name":"updated-policy","max_attempts":10,"status":"active"}`)
	defer close()
	res, err := newTestClient(t, srv).Retries.UpdatePolicy(context.Background(), 1, &edgecron.UpdateRetryPolicyRequest{MaxAttempts: int32Ptr(10)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MaxAttempts != 10 {
		t.Fatalf("expected max_attempts 10, got %d", res.MaxAttempts)
	}
}

func TestRetries_DeletePolicy(t *testing.T) {
	srv, close := mockServer(t, "/v1/retries/policies/1", http.MethodDelete, "")
	defer close()
	err := newTestClient(t, srv).Retries.DeletePolicy(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRetries_ListPolicies(t *testing.T) {
	srv, close := mockServer(t, "/v1/retries/policies", http.MethodGet, `{"total":0,"list":[]}`)
	defer close()
	_, err := newTestClient(t, srv).Retries.ListPolicies(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRetries_ListJobs(t *testing.T) {
	srv, close := mockServer(t, "/v1/retries/jobs", http.MethodGet, `{"total":0,"list":[]}`)
	defer close()
	_, err := newTestClient(t, srv).Retries.ListJobs(context.Background(), 1, 20, "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRetries_CancelJob(t *testing.T) {
	srv, close := mockServer(t, "/v1/retries/jobs/1/cancel", http.MethodPost, "")
	defer close()
	err := newTestClient(t, srv).Retries.CancelJob(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Subscription ---

func TestSubscription_Quota(t *testing.T) {
	srv, close := mockServer(t, "/v1/subscription/quota", http.MethodGet, `{"plan_code":"free","billing_cycle":"monthly","quota":10000,"used":0,"remaining":10000,"exceeded":false,"current_period_start":1700000000,"current_period_end":1730000000,"usage_percent":0}`)
	defer close()
	res, err := newTestClient(t, srv).Subscription.Quota(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.PlanCode != "free" || res.Quota != 10000 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestSubscription_Usage(t *testing.T) {
	srv, close := mockServer(t, "/v1/subscription/usage", http.MethodGet, `{"period":"2026-06","total_events":42,"items":[{"event_type":"inbound","period":"2026-06","count":30},{"event_type":"cron","period":"2026-06","count":12}]}`)
	defer close()
	res, err := newTestClient(t, srv).Subscription.Usage(context.Background(), "2026-06")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotalEvents != 42 || len(res.Items) != 2 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestSubscription_ResourceLimits(t *testing.T) {
	srv, close := mockServer(t, "/v1/subscription/resource-limits", http.MethodGet, `{"max_cron_jobs":5,"current_cron_jobs":2,"max_endpoints":-1,"current_endpoints":1,"log_retention_days":3}`)
	defer close()
	res, err := newTestClient(t, srv).Subscription.ResourceLimits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MaxCronJobs != 5 || res.CurrentCronJobs != 2 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

// --- Signing + Error Handling ---

func TestSigningSmoke(t *testing.T) {
	srv, close := mockServer(t, "/v1/retries/policies", http.MethodGet, `{"total":0,"list":[]}`)
	defer close()
	client := newTestClient(t, srv)
	_, err := client.Retries.ListPolicies(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":1001,"message":"bad request","request_id":"rid-001","data":null}`))
	}))
	defer srv.Close()
	client := edgecron.New("ak_3f9a2b1c7d4e8f0a", "test-secret", edgecron.WithBaseURL(srv.URL))
	_, err := client.Schedules.List(context.Background(), 1, 20, "")
	if err == nil {
		t.Fatal("expected APIError, got nil")
	}
	if apiErr, ok := edgecron.IsAPIError(err); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	} else if apiErr.Code != 1001 {
		t.Fatalf("expected code 1001, got %d", apiErr.Code)
	}
}

func TestKeyIDRejection(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid keyID")
		}
	}()
	edgecron.New("bad-key-id", "secret")
}

func TestOptions(t *testing.T) {
	c := edgecron.New("ak_3f9a2b1c7d4e8f0a", "secret", edgecron.WithTimeout(10*time.Second))
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestUserAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if !strings.HasPrefix(ua, "edgecron-go/") {
			t.Errorf("unexpected User-Agent: %s", ua)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0,"message":"success","request_id":"rid","data":null}`))
	}))
	defer srv.Close()
	client := edgecron.New("ak_3f9a2b1c7d4e8f0a", "secret", edgecron.WithBaseURL(srv.URL))
	client.Schedules.Pause(context.Background(), 1)
}

func strPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32 { return &i }
