package edgecron

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type RetriesService struct{ c *Client }

func (s *RetriesService) CreatePolicy(ctx context.Context, req *CreateRetryPolicyRequest) (*RetryPolicy, error) {
	var out RetryPolicy
	if err := s.c.do(ctx, http.MethodPost, "/v1/retries/policies", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *RetriesService) GetPolicy(ctx context.Context, id int64) (*RetryPolicy, error) {
	var out RetryPolicy
	if err := s.c.do(ctx, http.MethodGet, "/v1/retries/policies/"+strconv.FormatInt(id, 10), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *RetriesService) UpdatePolicy(ctx context.Context, id int64, req *UpdateRetryPolicyRequest) (*RetryPolicy, error) {
	var out RetryPolicy
	if err := s.c.do(ctx, http.MethodPatch, "/v1/retries/policies/"+strconv.FormatInt(id, 10), nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *RetriesService) DeletePolicy(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodDelete, "/v1/retries/policies/"+strconv.FormatInt(id, 10), nil, nil, nil)
}

func (s *RetriesService) ListPolicies(ctx context.Context) (*RetryPolicyList, error) {
	var out RetryPolicyList
	if err := s.c.do(ctx, http.MethodGet, "/v1/retries/policies", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *RetriesService) ListJobs(ctx context.Context, page, pageSize int, status string, deliveryID int64) (*RetryJobList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	if status != "" {
		q.Set("status", status)
	}
	if deliveryID > 0 {
		q.Set("delivery_id", strconv.FormatInt(deliveryID, 10))
	}
	var out RetryJobList
	if err := s.c.do(ctx, http.MethodGet, "/v1/retries/jobs", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *RetriesService) CancelJob(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodPost, "/v1/retries/jobs/"+strconv.FormatInt(id, 10)+"/cancel", nil, nil, nil)
}
