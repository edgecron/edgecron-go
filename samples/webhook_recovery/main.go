package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edgecron/edgecron-go"
)

func main() {
	keyID := os.Getenv("EDGECRON_KEY_ID")
	secret := os.Getenv("EDGECRON_SECRET")
	if keyID == "" || secret == "" {
		log.Fatal("EDGECRON_KEY_ID and EDGECRON_SECRET required")
	}

	client := edgecron.New(keyID, secret, edgecron.WithBaseURL("http://localhost:8888"))
	ctx := context.Background()

	policy, err := client.Retries.CreatePolicy(ctx, &edgecron.CreateRetryPolicyRequest{
		Name:            fmt.Sprintf("webhook-retry-%d", time.Now().Unix()),
		MaxAttempts:     5,
		BackoffType:     "exponential",
		InitialDelaySec: 10,
		MaxDelaySec:     3600,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retry policy created: %d (%s)\n", policy.ID, policy.Name)

	endpoint, err := client.Endpoints.Create(ctx, &edgecron.CreateEndpointRequest{
		Name:          "recovery-webhook",
		URL:           "https://httpbin.org/post",
		RetryPolicyID: policy.ID,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Endpoint created: %d (retry_policy_id=%d)\n", endpoint.ID, endpoint.RetryPolicyID)

	deliveries, err := client.Deliveries.List(ctx, 1, 10, "", 0, endpoint.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total deliveries for endpoint: %d\n", deliveries.Total)

	failedDeliveries, err := client.Deliveries.List(ctx, 1, 5, "failed", 0, endpoint.ID)
	if err != nil {
		log.Fatal(err)
	}
	if len(failedDeliveries.List) == 0 {
		fmt.Println("No failed deliveries found — retry skipped")
		return
	}
	failed := failedDeliveries.List[0]
	result, err := client.Deliveries.Retry(ctx, failed.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retry scheduled: delivery_id=%d retry_job_id=%d status=%s\n",
		result.DeliveryID, result.RetryJobID, result.Status)
}
