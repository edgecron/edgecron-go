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

	schedule, err := client.Schedules.Create(ctx, &edgecron.CreateScheduleRequest{
		Name:     fmt.Sprintf("my-cron-%d", time.Now().Unix()),
		CronExpr: "*/5 * * * *",
		Timezone: "UTC",
		Payload:  `{"task": "sync"}`,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Schedule created: %d (%s)\n", schedule.ID, schedule.Name)

	schedule, err = client.Schedules.Get(ctx, schedule.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Schedule: %s | cron: %s | status: %s\n", schedule.Name, schedule.CronExpr, schedule.Status)

	list, err := client.Schedules.List(ctx, 1, 10, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total schedules: %d\n", list.Total)

	updatedName := schedule.Name + "-updated"
	updatedSchedule, err := client.Schedules.Update(ctx, schedule.ID, &edgecron.UpdateScheduleRequest{
		Name: &updatedName,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Schedule updated: %s\n", updatedSchedule.Name)

	if err := client.Schedules.Pause(ctx, schedule.ID); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Schedule paused")

	if err := client.Schedules.Resume(ctx, schedule.ID); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Schedule resumed")

	if err := client.Schedules.Delete(ctx, schedule.ID); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Schedule deleted")
}
