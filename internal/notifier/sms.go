package notifier

import (
	"fmt"
	"log"
	"os"

	models "reminder/internal/model"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSMS(to string, reminder models.Reminder) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER")
	if accountSid == "" || authToken == "" || fromPhone == "" {
		return fmt.Errorf("Twilio environment variables are not set")
	}
    
	client := twilio.NewRestClient()
	params := &api.CreateMessageParams{}
	params.SetBody(reminder.Message)
	params.SetFrom(fromPhone)
	params.SetTo(to)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	} else {
		if resp.Body != nil {
			log.Printf("SMS sent successfully: %s", *resp.Body)
		} else {
			log.Println("SMS sent successfully")
		}
	}
	return nil
}