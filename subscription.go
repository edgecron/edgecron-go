package edgecron

import (
	"context"
	"net/http"
	"net/url"
)

// SubscriptionService handles subscription quota, usage, and resource limit queries.
type SubscriptionService struct{ c *Client }

// Quota returns the current subscription quota and usage for the authenticated app.
func (s *SubscriptionService) Quota(ctx context.Context) (*SubscriptionQuota, error) {
	var out SubscriptionQuota
	if err := s.c.do(ctx, http.MethodGet, "/v1/subscription/quota", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Usage returns event usage records for a given period (YYYY-MM). Empty period defaults to current month.
func (s *SubscriptionService) Usage(ctx context.Context, period string) (*UsageRecords, error) {
	q := url.Values{}
	if period != "" {
		q.Set("period", period)
	}
	var out UsageRecords
	if err := s.c.do(ctx, http.MethodGet, "/v1/subscription/usage", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResourceLimits returns the resource limits and current usage for the authenticated app.
func (s *SubscriptionService) ResourceLimits(ctx context.Context) (*ResourceLimits, error) {
	var out ResourceLimits
	if err := s.c.do(ctx, http.MethodGet, "/v1/subscription/resource-limits", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
