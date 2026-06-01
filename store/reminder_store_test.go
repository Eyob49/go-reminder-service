package store

import (
	"database/sql"
	"os"
	models "reminder/internal/model"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

// Tests for InsertReminder
func SetupTestDB(t *testing.T) *sql.DB {
	godotenv.Load("../.env")
	// Setup in-memory database for all tests
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	err = InitDB(db)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	_, err = db.Exec(`DELETE FROM reminders`)
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}
	return db
}

func TestInsertReminder(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()
	reminder := models.Reminder{
		Title: "Test Reminder",
		Message: "This is a test reminder",
		Email: "test@example.com",
		Phone: "1234567890",
		Channel: "email",
		SendAt: time.Now().Add(1 * time.Hour),
	}
	err := InsertReminder(db, reminder)
	if err != nil {
		t.Fatalf("Failed to insert reminder: %v", err)
	}

}

// Test for GetReminders

func TestGetReminders(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()
	reminder := models.Reminder{
		Title: "Test Reminder",
		Message: "This is a test Reminder",
		Email: "test@example.com",
		Phone: "1234567890",
		Channel: "email",
		SendAt: time.Now().Add(1 * time.Hour),
	}
	err := InsertReminder(db, reminder)
	if err != nil {
		t.Fatalf("Failed to insert reminder: %v", err)
	}
	reminders, err := GetReminders(db)
	if err != nil {
		t.Fatalf("Failed to get reminders: %v", err)
	}
	if reminders[0].Title != reminder.Title {
		t.Fatalf("Expected reminder with title '%s', but got '%s'", reminder.Title, reminders[0].Title)
	}

}

// Test for DeleteReminder
func TestDeleteReminder(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()
	reminder := models.Reminder{
		Title: "Test Reminder",
		Message: "This is a test Reminder",
		Email: "test@example.com",
		Phone: "1234567890",
		Channel: "email",
		SendAt: time.Now().Add(1 * time.Hour),
	}
	err := InsertReminder(db, reminder)
	if err != nil {
		t.Fatalf("Failed to insert reminder: %v", err)
	}
	var id int
	err = db.QueryRow(`SELECT id FROM reminders WHERE title = $1`, reminder.Title).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to query reminder ID: %v", err)
	}
	err = DeleteReminder(db, id)
	if err != nil {
		t.Fatalf("Failed to delete reminder: %v", err)
	}
	reminders, err := GetReminders(db)
	if err != nil {
		t.Fatalf("Failed to get reminders: %v", err)
	}
	if len(reminders) != 0 {
		t.Fatalf("Expected 0 reminders, but got %d", len(reminders))
	}
}

// Test for DeleteReminder with non-existent ID
func TestDeleteNonExistentReminder(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()
	err := DeleteReminder(db, 999)
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound, but got %v", err)
	}
}
