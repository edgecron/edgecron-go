package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	endpoint, err := client.Endpoints.Create(ctx, &edgecron.CreateEndpointRequest{
		Name: "my-webhook",
		URL:  "https://httpbin.org/post",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Endpoint created: %d\n", endpoint.ID)

	task, err := client.Tasks.Create(ctx, &edgecron.CreateTaskRequest{
		EndpointID: endpoint.ID,
		Payload:    `{"order_id": "ord_9102"}`,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Task created: %d\n", task.ID)
}
