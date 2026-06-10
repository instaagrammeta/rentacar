package service

import (
	"fmt"
	"log"
	"time"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// SendClientSMS sends a custom SMS to a client (manual, by an operator/admin).
func (s *Service) SendClientSMS(clientID uint, message string, actor *Actor) error {
	if message == "" {
		return apperr.Validation("Текст сообщения обязателен")
	}
	client, err := s.GetClient(clientID)
	if err != nil {
		return err
	}
	if client.Phone == "" {
		return apperr.Validation("У клиента не указан номер телефона")
	}
	if err := s.SMS.Send(client.Phone, message); err != nil {
		return apperr.New(err.Error(), 502)
	}
	s.record(actor, "send_sms", "client", uintPtr(client.ID), "SMS: "+message)
	return nil
}

// ProcessRentalReminders sends automatic reminders to clients whose active
// rental is about to expire: one ~1 hour before, one ~30 minutes before.
// It is idempotent thanks to the reminder flags on the rental.
func (s *Service) ProcessRentalReminders() {
	if s.SMS == nil || !s.SMS.Configured() {
		return
	}
	now := time.Now()

	var rentals []models.Rental
	// Candidates: active rentals with a due time within the next hour.
	if err := s.DB.Preload("Client").Preload("Car").
		Where("status = ? AND due_at IS NOT NULL AND due_at <= ?", models.RentalActive, now.Add(time.Hour)).
		Find(&rentals).Error; err != nil {
		log.Printf("reminders: query failed: %v", err)
		return
	}

	for i := range rentals {
		r := &rentals[i]
		if r.DueAt == nil || r.Client == nil || r.Client.Phone == "" {
			continue
		}
		remaining := r.DueAt.Sub(now)
		carName := ""
		if r.Car != nil {
			carName = r.Car.DisplayName()
		}
		dueStr := r.DueAt.Format("02.01.2006 15:04")

		// 1-hour reminder.
		if !r.Reminder1hSent && remaining <= time.Hour {
			msg := fmt.Sprintf("Уважаемый %s! Срок аренды автомобиля %s истекает в %s (через ~1 час). Просим вернуть автомобиль вовремя.",
				r.Client.FirstName, carName, dueStr)
			if err := s.SMS.Send(r.Client.Phone, msg); err != nil {
				log.Printf("reminders: 1h SMS failed for rental %d: %v", r.ID, err)
			} else {
				s.DB.Model(&models.Rental{}).Where("id = ?", r.ID).Update("reminder_1h_sent", true)
			}
		}

		// 30-minute reminder.
		if !r.Reminder30mSent && remaining <= 30*time.Minute {
			msg := fmt.Sprintf("Уважаемый %s! До окончания аренды автомобиля %s осталось около 30 минут (до %s). Пожалуйста, верните автомобиль.",
				r.Client.FirstName, carName, dueStr)
			if err := s.SMS.Send(r.Client.Phone, msg); err != nil {
				log.Printf("reminders: 30m SMS failed for rental %d: %v", r.ID, err)
			} else {
				s.DB.Model(&models.Rental{}).Where("id = ?", r.ID).Update("reminder_30m_sent", true)
			}
		}
	}
}
