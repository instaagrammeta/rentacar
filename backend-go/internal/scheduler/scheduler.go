// Package scheduler runs background jobs such as the daily database backup.
package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"

	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// Start launches the cron scheduler. It returns the cron instance so the caller
// can stop it on shutdown. Returns nil when scheduled backups are disabled.
func Start(svc *service.Service) *cron.Cron {
	if !svc.Cfg.EnableScheduledBackups {
		return nil
	}
	c := cron.New()
	// Every day at 02:00.
	_, err := c.AddFunc("0 2 * * *", func() {
		if path, err := svc.CreateBackup(); err != nil {
			log.Printf("scheduled backup failed: %v", err)
		} else {
			log.Printf("scheduled backup created: %s", path)
		}
	})
	if err != nil {
		log.Printf("could not schedule backup job: %v", err)
		return nil
	}
	c.Start()
	log.Println("Планировщик резервного копирования запущен (ежедневно в 02:00)")
	return c
}
