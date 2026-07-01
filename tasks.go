package edgecron

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type TasksService struct{ c *Client }

func (s *TasksService) Create(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
	var out Task
	if err := s.c.do(ctx, http.MethodPost, "/v1/tasks", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TasksService) Get(ctx context.Context, id int64) (*Task, error) {
	var out Task
	if err := s.c.do(ctx, http.MethodGet, "/v1/tasks/"+strconv.FormatInt(id, 10), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TasksService) List(ctx context.Context, page, pageSize int, status string, scheduleID, eventID int64) (*TaskList, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	if status != "" {
		q.Set("status", status)
	}
	if scheduleID > 0 {
		q.Set("schedule_id", strconv.FormatInt(scheduleID, 10))
	}
	if eventID > 0 {
		q.Set("event_id", strconv.FormatInt(eventID, 10))
	}
	var out TaskList
	if err := s.c.do(ctx, http.MethodGet, "/v1/tasks", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TasksService) Cancel(ctx context.Context, id int64) error {
	return s.c.do(ctx, http.MethodPost, "/v1/tasks/"+strconv.FormatInt(id, 10)+"/cancel", nil, nil, nil)
}
