package notifier

import (
	"fmt"
	"log"
	"os"
	"reminder/internal/models"

	"github.com/resend/resend-go/v3"
)

func SendEmail(to string, reminder models.Reminder) error {
    apiKey := os.Getenv("RESEND_API_KEY")
    if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY environment variable is not set")
	}
    client := resend.NewClient(apiKey)

    params := &resend.SendEmailRequest{
        From:    "onboarding@resend.dev",
        To:      []string{reminder.Email},
        Subject: reminder.Title,
        Html:    "<p>" + reminder.Message + "</p>",
    }

    sent, err := client.Emails.Send(params)
    if err != nil {
        return fmt.Errorf("failed to send email: %v", err)
    }
    log.Printf("email sent: %s", sent.Id)
    return nil
}