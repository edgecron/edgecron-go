package edgecron

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type EndpointsService struct{ c *Client }

func (s *EndpointsService) Create(ctx context.Context, req *CreateEndpointRequest) (*WebhookEndpoint, error) {
	var out WebhookEndpoint
	if err := s.c.do(ctx, http.MethodPost, "/v1/endpoints", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *EndpointsService) Get(ctx context.Context, id int64) (*WebhookEndpoint, error) {
	var out WebhookEndpoint
	if err := s.c.do(ctx, http.MethodGet, "/v1/endpoints/"+strconv.FormatInt(id, 10), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *EndpointsService) Update(ctx context.Context, id int64, req *UpdateEndpointRequest) (*WebhookEndpoint, error) {
	var out WebhookEndpoint
	if err := s.c.do(ctx, http.MethodPatch, "/v1/endpoints/"+strconv.FormatInt(id, 10), nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *EndpointsService) List(ctx context.Context, page, pageSize int, status string) (*EndpointList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	if status != "" {
		q.Set("status", status)
	}
	var out EndpointList
	if err := s.c.do(ctx, http.MethodGet, "/v1/endpoints", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *EndpointsService) Delete(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodDelete, "/v1/endpoints/"+strconv.FormatInt(id, 10), nil, nil, nil)
}

func (s *EndpointsService) Enable(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodPost, "/v1/endpoints/"+strconv.FormatInt(id, 10)+"/enable", nil, nil, nil)
}

func (s *EndpointsService) Disable(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodPost, "/v1/endpoints/"+strconv.FormatInt(id, 10)+"/disable", nil, nil, nil)
}
