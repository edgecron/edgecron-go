# EdgeCron Go SDK

Official Go SDK for the EdgeCron webhook scheduling and callback delivery platform.

Schedule delayed HTTP requests, deliver webhooks reliably, and automatically retry failed calls — with full execution history so nothing gets lost.

## Install

```bash
go get github.com/edgecron/edgecron-go@latest
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/edgecron/edgecron-go"
)

func main() {
    client := edgecron.New("ak_xxx", "sk_xxx")

    schedule, err := client.Schedules.Create(context.Background(), &edgecron.CreateScheduleRequest{
        Name:     "my-schedule",
        CronExpr: "*/5 * * * *",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Schedule created: %d\n", schedule.ID)
}
```

## Modules

| Client method         | Description                        |
| --------------------- | ---------------------------------- |
| `client.Schedules.*`  | Cron schedule CRUD, pause, resume  |
| `client.Tasks.*`      | Task execution instances, cancel   |
| `client.Events.*`     | Event publishing and management    |
| `client.Endpoints.*`  | Webhook endpoint configuration     |
| `client.Deliveries.*` | Delivery attempt records and retry |
| `client.Retries.*`    | Retry policies and jobs            |
| `client.Subscription.*` | Quota, usage, and resource limits |

## Configuration

- `edgecron.WithBaseURL(url)` — override API base URL for private deployments
- `edgecron.WithTimeout(duration)` — set HTTP client timeout
- `edgecron.WithHTTPClient(client)` — provide custom `*http.Client`

## Error Handling

Service-side business errors are returned as `*edgecron.APIError`.

```go
if apiErr, ok := edgecron.IsAPIError(err); ok {
    fmt.Println(apiErr.Code, apiErr.Message, apiErr.RequestID)
}
```

## Security Notice

This is a server-side SDK. Never expose `secret` in browsers, mobile apps, or other untrusted clients.
