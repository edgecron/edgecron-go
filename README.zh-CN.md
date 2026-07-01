# EdgeCron Go SDK

EdgeCron Go SDK 是 EdgeCron Webhook 调度与回调投递平台的官方 Go 客户端。

调度延迟 HTTP 请求，可靠投递 Webhook，自动重试失败调用 — 完整执行历史，确保不遗漏。

## 安装

```bash
go get github.com/edgecron/edgecron-go@latest
```

## 快速开始

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

## 模块说明

| 客户端方法            | 说明                         |
| --------------------- | ---------------------------- |
| `client.Schedules.*`  | Cron 调度器 CRUD、暂停、恢复 |
| `client.Tasks.*`      | 任务执行实例、取消           |
| `client.Events.*`     | 事件发布与管理               |
| `client.Endpoints.*`  | Webhook 端点配置             |
| `client.Deliveries.*` | 投递记录与手动重试           |
| `client.Retries.*`    | 重试策略与任务               |
| `client.Subscription.*` | 配额、用量与资源限制           |

## 配置项

- `edgecron.WithBaseURL(url)` — 覆盖 API 地址，适合私有部署
- `edgecron.WithTimeout(duration)` — 覆盖超时时间
- `edgecron.WithHTTPClient(client)` — 传入自定义 `http.Client`

## 错误处理

服务端业务错误会返回 `*edgecron.APIError`。

```go
if apiErr, ok := edgecron.IsAPIError(err); ok {
    fmt.Println(apiErr.Code, apiErr.Message, apiErr.RequestID)
}
```

## 安全说明

这是服务端 SDK，不要在浏览器、小程序、移动端或其他不可信客户端暴露 `secret`。
