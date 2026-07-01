package edgecron

// --- Schedules ---

type Schedule struct {
	ID            int64            `json:"id"`
	AppID         string           `json:"app_id"`
	Name          string           `json:"name"`
	CronExpr      string           `json:"cron_expr"`
	Timezone      string           `json:"timezone"`
	Payload       string           `json:"payload"`
	Status        string           `json:"status"`
	NextRunAt     int64            `json:"next_run_at"`
	EndpointIDs   []int64          `json:"endpoint_ids"`
	EndpointNames map[int64]string `json:"endpoint_names"`
	CreatedAt     int64            `json:"created_at"`
	UpdatedAt     int64            `json:"updated_at"`
}

type CreateScheduleRequest struct {
	Name        string  `json:"name"`
	CronExpr    string  `json:"cron_expr"`
	Timezone    string  `json:"timezone,omitempty"`
	Payload     string  `json:"payload,omitempty"`
	EndpointIDs []int64 `json:"endpoint_ids,omitempty"`
}

type UpdateScheduleRequest struct {
	Name        *string  `json:"name,omitempty"`
	CronExpr    *string  `json:"cron_expr,omitempty"`
	Timezone    *string  `json:"timezone,omitempty"`
	Payload     *string  `json:"payload,omitempty"`
	EndpointIDs *[]int64 `json:"endpoint_ids,omitempty"`
}

type ScheduleList struct {
	Total int64      `json:"total"`
	List  []Schedule `json:"list"`
}

// --- Tasks ---

type Task struct {
	ID         int64  `json:"id"`
	AppID      string `json:"app_id"`
	ScheduleID int64  `json:"schedule_id"`
	EventID    int64  `json:"event_id"`
	EndpointID int64  `json:"endpoint_id"`
	TaskType   string `json:"task_type"`
	Payload    string `json:"payload"`
	Status     string `json:"status"`
	RunAt      int64  `json:"run_at"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type CreateTaskRequest struct {
	EndpointID int64  `json:"endpoint_id"`
	Payload    string `json:"payload,omitempty"`
	RunAt      int64  `json:"run_at,omitempty"`
}

type TaskList struct {
	Total int64  `json:"total"`
	List  []Task `json:"list"`
}

// --- Events ---

type Event struct {
	ID        int64  `json:"id"`
	AppID     string `json:"app_id"`
	EventName string `json:"event_name"`
	EventKey  string `json:"event_key"`
	Payload   string `json:"payload"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type PublishEventResult struct {
	ID          int64  `json:"id"`
	AppID       string `json:"app_id"`
	EventName   string `json:"event_name"`
	EventKey    string `json:"event_key"`
	Payload     string `json:"payload"`
	Status      string `json:"status"`
	FanoutCount int    `json:"fanout_count"`
	CreatedAt   int64  `json:"created_at"`
}

type PublishEventRequest struct {
	EventName string `json:"event_name"`
	EventKey  string `json:"event_key"`
	Payload   string `json:"payload,omitempty"`
}

type EventList struct {
	Total int64   `json:"total"`
	List  []Event `json:"list"`
}

// --- Endpoints ---

type WebhookEndpoint struct {
	ID            int64  `json:"id"`
	AppID         string `json:"app_id"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	Method        string `json:"method"`
	Headers       string `json:"headers"`
	Secret        string `json:"secret"`
	TimeoutMs     int32  `json:"timeout_ms"`
	RetryPolicyID int64  `json:"retry_policy_id"`
	FilterEvents  string `json:"filter_events"`
	Status        string `json:"status"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

type CreateEndpointRequest struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Method        string `json:"method,omitempty"`
	Headers       string `json:"headers,omitempty"`
	Secret        string `json:"secret,omitempty"`
	TimeoutMs     int32  `json:"timeout_ms,omitempty"`
	RetryPolicyID int64  `json:"retry_policy_id,omitempty"`
	FilterEvents  string `json:"filter_events,omitempty"`
}

type UpdateEndpointRequest struct {
	Name          *string `json:"name,omitempty"`
	URL           *string `json:"url,omitempty"`
	Method        *string `json:"method,omitempty"`
	Headers       *string `json:"headers,omitempty"`
	Secret        *string `json:"secret,omitempty"`
	TimeoutMs     *int32  `json:"timeout_ms,omitempty"`
	RetryPolicyID *int64  `json:"retry_policy_id,omitempty"`
	FilterEvents  *string `json:"filter_events,omitempty"`
}

type EndpointList struct {
	Total int64             `json:"total"`
	List  []WebhookEndpoint `json:"list"`
}

// --- Deliveries ---

type Delivery struct {
	ID              int64  `json:"id"`
	AppID           string `json:"app_id"`
	TaskID          int64  `json:"task_id"`
	EndpointID      int64  `json:"endpoint_id"`
	Attempt         int32  `json:"attempt"`
	Status          string `json:"status"`
	HTTPStatus      int32  `json:"http_status"`
	LatencyMs       int32  `json:"latency_ms"`
	RequestBodyHash string `json:"request_body_hash"`
	ErrorMessage    string `json:"error_message"`
	NextRetryAt     int64  `json:"next_retry_at"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

type DeliveryList struct {
	Total int64      `json:"total"`
	List  []Delivery `json:"list"`
}

type RetryDeliveryResult struct {
	DeliveryID int64  `json:"delivery_id"`
	RetryJobID int64  `json:"retry_job_id"`
	Status     string `json:"status"`
}

// --- Retries ---

type RetryPolicy struct {
	ID              int64  `json:"id"`
	AppID           string `json:"app_id"`
	Name            string `json:"name"`
	MaxAttempts     int32  `json:"max_attempts"`
	BackoffType     string `json:"backoff_type"`
	InitialDelaySec int32  `json:"initial_delay_sec"`
	MaxDelaySec     int32  `json:"max_delay_sec"`
	Status          string `json:"status"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

type CreateRetryPolicyRequest struct {
	Name            string `json:"name"`
	MaxAttempts     int32  `json:"max_attempts,omitempty"`
	BackoffType     string `json:"backoff_type,omitempty"`
	InitialDelaySec int32  `json:"initial_delay_sec,omitempty"`
	MaxDelaySec     int32  `json:"max_delay_sec,omitempty"`
}

type UpdateRetryPolicyRequest struct {
	Name            *string `json:"name,omitempty"`
	MaxAttempts     *int32  `json:"max_attempts,omitempty"`
	BackoffType     *string `json:"backoff_type,omitempty"`
	InitialDelaySec *int32  `json:"initial_delay_sec,omitempty"`
	MaxDelaySec     *int32  `json:"max_delay_sec,omitempty"`
	Status          *string `json:"status,omitempty"`
}

type RetryPolicyList struct {
	Total int64         `json:"total"`
	List  []RetryPolicy `json:"list"`
}

type RetryJob struct {
	ID          int64  `json:"id"`
	AppID       string `json:"app_id"`
	DeliveryID  int64  `json:"delivery_id"`
	Attempt     int32  `json:"attempt"`
	Status      string `json:"status"`
	RunAt       int64  `json:"run_at"`
	LockedUntil int64  `json:"locked_until"`
	LastError   string `json:"last_error"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type RetryJobList struct {
	Total int64      `json:"total"`
	List  []RetryJob `json:"list"`
}

// --- Subscription ---

// SubscriptionQuota represents the current subscription quota and usage.
type SubscriptionQuota struct {
	PlanCode           string  `json:"plan_code"`
	BillingCycle       string  `json:"billing_cycle"`
	Quota              int64   `json:"quota"`
	Used               int64   `json:"used"`
	Remaining          int64   `json:"remaining"`
	Exceeded           bool    `json:"exceeded"`
	CurrentPeriodStart int64   `json:"current_period_start"`
	CurrentPeriodEnd   int64   `json:"current_period_end"`
	UsagePercent       float64 `json:"usage_percent"`
}

// UsageRecordItem is a single event-type usage count for a given period.
type UsageRecordItem struct {
	EventType string `json:"event_type"`
	Period    string `json:"period"`
	Count     int64  `json:"count"`
}

// UsageRecords contains event usage records for a billing period.
type UsageRecords struct {
	Period      string            `json:"period"`
	TotalEvents int64             `json:"total_events"`
	Items       []UsageRecordItem `json:"items"`
}

// ResourceLimits represents the subscription plan's resource limits and current usage.
type ResourceLimits struct {
	MaxCronJobs      int `json:"max_cron_jobs"`
	CurrentCronJobs  int `json:"current_cron_jobs"`
	MaxEndpoints     int `json:"max_endpoints"`
	CurrentEndpoints int `json:"current_endpoints"`
	LogRetentionDays int `json:"log_retention_days"`
}
