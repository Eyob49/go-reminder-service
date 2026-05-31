package scheduler

import (
	"context"
	"database/sql"
	"log"
	"reminder/internal/notifier"
	"reminder/store"
	"time"
)

func StartScheduler(ctx context.Context, db *sql.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler shutting down")
			return
		case <-ticker.C:
			log.Println("Checking for due reminders...")
			reminders, err := store.GetDueReminders(db)
			if err != nil {
				log.Printf("Failed to retrieve due reminders: %v", err)
				continue
			}
			for _, reminder := range reminders {
				switch reminder.Channel {
				case "email":
					err = notifier.SendEmail(reminder.Email, reminder)
					if err != nil {
						log.Printf("Failed to send email for reminder ID %d: %v", reminder.ID, err)
					}
				case "sms":
					err = notifier.SendSMS(reminder.Phone, reminder)
					if err != nil {
						log.Printf("Failed to send SMS for reminder ID %d: %v", reminder.ID, err)
					}
				case "both":
					err = notifier.SendEmail(reminder.Email, reminder)
					if err != nil {
						log.Printf("Failed to send email for reminder ID %d: %v", reminder.ID, err)
					}
					err = notifier.SendSMS(reminder.Phone, reminder)
					if err != nil {
						log.Printf("Failed to send SMS for reminder ID %d: %v", reminder.ID, err)
					}
				default:
					log.Printf("Unknown channel for reminder ID %d: %s", reminder.ID, reminder.Channel)
				}
				err = store.MarkReminderAsSent(db, reminder.ID)
				if err != nil {
					log.Printf("Error marking reminder as sent: %v", err)
				}
			}
		}
	}
}
		    
