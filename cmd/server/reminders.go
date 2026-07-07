package main

import (
	"context"
	"log/slog"
	"time"

	sv "github.com/dariojcalo91/gym-backend-go-ver/internal/service"
)

// runDailyReminders corre una vez al arrancar (útil en Render free tier, donde
// el proceso se duerme y despierta) y luego cada 24h mientras el proceso viva.
func runDailyReminders(ctx context.Context, reminders *sv.ReminderService) {
	run := func() {
		if err := reminders.RunDaily(ctx); err != nil {
			slog.Error("daily reminder job failed", "error", err)
		}
	}

	run()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			run()
		case <-ctx.Done():
			return
		}
	}
}
