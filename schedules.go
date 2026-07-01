package edgecron

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type SchedulesService struct{ c *Client }

func (s *SchedulesService) Create(ctx context.Context, req *CreateScheduleRequest) (*Schedule, error) {
	var out Schedule
	if err := s.c.do(ctx, http.MethodPost, "/v1/schedules", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *SchedulesService) Update(ctx context.Context, id int64, req *UpdateScheduleRequest) (*Schedule, error) {
	var out Schedule
	if err := s.c.do(ctx, http.MethodPatch, "/v1/schedules/"+strconv.FormatInt(id, 10), nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *SchedulesService) Get(ctx context.Context, id int64) (*Schedule, error) {
	var out Schedule
	if err := s.c.do(ctx, http.MethodGet, "/v1/schedules/"+strconv.FormatInt(id, 10), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *SchedulesService) List(ctx context.Context, page, pageSize int, status string) (*ScheduleList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	if status != "" {
		q.Set("status", status)
	}
	var out ScheduleList
	if err := s.c.do(ctx, http.MethodGet, "/v1/schedules", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *SchedulesService) Delete(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodDelete, "/v1/schedules/"+strconv.FormatInt(id, 10), nil, nil, nil)
}

func (s *SchedulesService) Pause(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodPost, "/v1/schedules/"+strconv.FormatInt(id, 10)+"/pause", nil, nil, nil)
}

func (s *SchedulesService) Resume(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodPost, "/v1/schedules/"+strconv.FormatInt(id, 10)+"/resume", nil, nil, nil)
}
