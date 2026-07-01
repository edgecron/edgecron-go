package edgecron

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type DeliveriesService struct{ c *Client }

func (s *DeliveriesService) List(ctx context.Context, page, pageSize int, status string, taskID, endpointID int64) (*DeliveryList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	if status != "" {
		q.Set("status", status)
	}
	if taskID > 0 {
		q.Set("task_id", strconv.FormatInt(taskID, 10))
	}
	if endpointID > 0 {
		q.Set("endpoint_id", strconv.FormatInt(endpointID, 10))
	}
	var out DeliveryList
	if err := s.c.do(ctx, http.MethodGet, "/v1/deliveries", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *DeliveriesService) Retry(ctx context.Context, id int64) (*RetryDeliveryResult, error) {
	var out RetryDeliveryResult
	if err := s.c.do(ctx, http.MethodPost, "/v1/deliveries/"+strconv.FormatInt(id, 10)+"/retry", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
