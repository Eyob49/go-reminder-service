package store

import (
	"database/sql"
	"fmt"
	"log"
	models "reminder/internal/model"
	"time"

	_ "modernc.org/sqlite"
)

var ErrNotFound = fmt.Errorf("reminder not found")

func InitDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS reminders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		message TEXT,
		email TEXT,
		phone TEXT,
		channel TEXT,
		send_at TEXT,
		sent INTEGER DEFAULT 0
		)`)
    if err != nil {
	return err
    }
	log.Println("Database intialized successfully")
	return nil
}

// Function to Insert a new reminder into the database
func InsertReminder(db *sql.DB, reminder models.Reminder) error {
	_, err := db.Exec(`INSERT INTO reminders (title, message, email, phone, channel, send_at) VALUES (?, ?, ?, ?, ?, ?)`, reminder.Title, reminder.Message, reminder.Email, reminder.Phone, reminder.Channel, reminder.SendAt.UTC().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	log.Println("Reminder inserted into database successfully")
	return nil
}

// Function to Retrieve all reminders from the database
func GetReminders(db *sql.DB) ([]models.Reminder, error) {

	rows, err := db.Query(`SELECT id, title, message, email, phone, channel, send_at, sent FROM reminders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []models.Reminder
	for rows.Next() {
		var reminder models.Reminder
		var sendAtStr string
		err := rows.Scan(&reminder.ID, &reminder.Title, &reminder.Message, &reminder.Email, &reminder.Phone, &reminder.Channel, &sendAtStr, &reminder.Sent)
		if err != nil {
			return nil, err
		}
		reminder.SendAt, err = time.Parse("2006-01-02 15:04:05", sendAtStr)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	return reminders, nil
}

// Function to delete a reminder from the database 
func DeleteReminder(db *sql.DB, id int) error {
	result, err := db.Exec(`DELETE FROM reminders WHERE id = ?`, id)
	if err != nil {
		return err
	}
    rowsAffected,_ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	} 
	return nil
}

// Function that returns only reminders where send_at is in the past and sent = 0
func GetDueReminders(db *sql.DB) ([]models.Reminder, error) {
	rows, err := db.Query(`SELECT id, title, message, email, phone, channel, send_at, sent FROM reminders WHERE send_at <= datetime('now') AND sent = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []models.Reminder
	for rows.Next() {
		var reminder models.Reminder
		var sendAtStr string
		err := rows.Scan(&reminder.ID, &reminder.Title, &reminder.Message, &reminder.Email, &reminder.Phone, &reminder.Channel, &sendAtStr, &reminder.Sent)
		if err != nil {
			return nil, err
		}
		reminder.SendAt, err = time.Parse("2006-01-02 15:04:05", sendAtStr)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	return reminders, nil
}

// Function to mark a reminder as sent in the database
func MarkReminderAsSent(db *sql.DB, id int) error {
	_, err := db.Exec(`UPDATE reminders SET sent = 1 WHERE id = ?`, id)
	return err
}
