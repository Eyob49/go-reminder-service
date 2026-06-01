package models

import "time"

type Reminder struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	SendAt  time.Time `json:"sendAt"`
	Sent    bool     `json:"sent"`
	Channel string `json:"channel"`
}