// Package scheduler runs background jobs such as the daily database backup.
package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"

	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// Start launches the cron scheduler with the daily backup job and the rental
// SMS-reminder job. It returns the cron instance so the caller can stop it on
// shutdown (or nil when there is nothing to schedule).
func Start(svc *service.Service) *cron.Cron {
	c := cron.New()
	scheduled := false

	if svc.Cfg.EnableScheduledBackups {
		if _, err := c.AddFunc("0 2 * * *", func() {
			if path, err := svc.CreateBackup(); err != nil {
				log.Printf("scheduled backup failed: %v", err)
			} else {
				log.Printf("scheduled backup created: %s", path)
			}
		}); err != nil {
			log.Printf("could not schedule backup job: %v", err)
		} else {
			scheduled = true
			log.Println("Планировщик резервного копирования запущен (ежедневно в 02:00)")
		}
	}

	if svc.SMS != nil && svc.SMS.Configured() {
		if _, err := c.AddFunc("@every 5m", func() {
			svc.ProcessRentalReminders()
		}); err != nil {
			log.Printf("could not schedule SMS reminder job: %v", err)
		} else {
			scheduled = true
			log.Println("Планировщик SMS-напоминаний запущен (каждые 5 минут)")
		}
	}

	if !scheduled {
		return nil
	}
	c.Start()
	return c
}
