package edgecron

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type EventsService struct{ c *Client }

func (s *EventsService) Publish(ctx context.Context, req *PublishEventRequest) (*PublishEventResult, error) {
	var out PublishEventResult
	if err := s.c.do(ctx, http.MethodPost, "/v1/events", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *EventsService) Get(ctx context.Context, id int64) (*Event, error) {
	var out Event
	if err := s.c.do(ctx, http.MethodGet, "/v1/events/"+strconv.FormatInt(id, 10), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *EventsService) List(ctx context.Context, page, pageSize int, eventName, status string) (*EventList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	if eventName != "" {
		q.Set("event_name", eventName)
	}
	if status != "" {
		q.Set("status", status)
	}
	var out EventList
	if err := s.c.do(ctx, http.MethodGet, "/v1/events", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *EventsService) Enable(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodPost, "/v1/events/"+strconv.FormatInt(id, 10)+"/enable", nil, nil, nil)
}

func (s *EventsService) Disable(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodPost, "/v1/events/"+strconv.FormatInt(id, 10)+"/disable", nil, nil, nil)
}

func (s *EventsService) Delete(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodDelete, "/v1/events/"+strconv.FormatInt(id, 10), nil, nil, nil)
}
