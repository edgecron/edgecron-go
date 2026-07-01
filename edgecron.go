// Package edgecron provides a Go SDK for the EdgeCron webhook platform.
//
// # Quick start
//
//	client := edgecron.New("ak_xxx", "sk_xxx")
//	schedule, err := client.Schedules.Create(ctx, &edgecron.CreateScheduleRequest{
//	    Name:     "my-schedule",
//	    CronExpr: "*/5 * * * *",
//	})
//
// # Authentication
//
// Every request is signed with HMAC-SHA256 using your secret key.
// The SDK handles signing automatically; you only need to supply keyID and secret.
//
// # Error handling
//
// Service errors are returned as *APIError with a numeric Code and Message.
// Use IsAPIError to inspect them:
//
//	if apiErr, ok := edgecron.IsAPIError(err); ok {
//	    fmt.Println(apiErr.Code, apiErr.Message)
//	}
//
// # Thread safety
//
// A Client is safe for concurrent use. Create one and reuse it across goroutines.
package edgecron

import (
	"net/http"
	"regexp"
	"time"
)

const (
	defaultBaseURL = "https://api.edgecron.com"
	Version        = "1.0.0"
)

var reKeyID = regexp.MustCompile(`^ak_[0-9a-zA-Z_]+$`)

type Client struct {
	keyID   string
	secret  string
	baseURL string
	http    *http.Client

	Schedules    *SchedulesService
	Tasks        *TasksService
	Events       *EventsService
	Endpoints    *EndpointsService
	Deliveries   *DeliveriesService
	Retries      *RetriesService
	Subscription *SubscriptionService
}

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.http = h }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

func New(keyID, secret string, opts ...Option) *Client {
	if !reKeyID.MatchString(keyID) {
		panic("edgecron: keyID must match ak_<hex>, got: " + keyID)
	}
	if secret == "" {
		panic("edgecron: secret must not be empty")
	}
	c := &Client{
		keyID:   keyID,
		secret:  secret,
		baseURL: defaultBaseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	c.Schedules = &SchedulesService{c: c}
	c.Tasks = &TasksService{c: c}
	c.Events = &EventsService{c: c}
	c.Endpoints = &EndpointsService{c: c}
	c.Deliveries = &DeliveriesService{c: c}
	c.Retries = &RetriesService{c: c}
	c.Subscription = &SubscriptionService{c: c}
	return c
}
