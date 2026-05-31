package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	models "reminder/internal/model"
	"reminder/store"
	"strconv"
	"strings"
	"time"
)

type ReminderHandler struct {
	DB *sql.DB
}

func  (h *ReminderHandler) HandleReminders(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case http.MethodGet:
		// Handle GET request to retrieve reminders
		reminders, err := store.GetReminders(h.DB)
		if err != nil {
			log.Printf("Failed to retrieve reminders: %v", err)
			http.Error(w, "Failed to retrieve reminders", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reminders)



	case http.MethodPost:
		// Handle POST request to create a new reminder
		var reminder models.Reminder
		err := json.NewDecoder(r.Body).Decode(&reminder)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		// Validate required fields
		if reminder.Title == "" || reminder.Message == "" || reminder.Channel == "" || reminder.SendAt.Equal((time.Time{})){
			http.Error(w, "Missing required fields", http.StatusUnprocessableEntity)
			return
		} 
		if reminder.Channel == "email" && reminder.Email == "" {
			http.Error(w, "Email is required for email channel", http.StatusUnprocessableEntity)
			return
		}
		if reminder.Channel == "sms" && reminder.Phone == "" {
			http.Error(w, "Phone number is required for sms channel", http.StatusUnprocessableEntity)
			return
		}
		log.Printf("SendAt: %v, Now: %v", reminder.SendAt, time.Now().UTC())
		if reminder.SendAt.Before(time.Now().UTC()) {
			http.Error(w, "SendAt must be a future time", http.StatusUnprocessableEntity)
			return
		}
		err = store.InsertReminder(h.DB, reminder)
		if err != nil {
			http.Error(w, "Failed to create reminder", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(reminder)
	case http.MethodDelete:
		// Handle DELETE request to delete a reminder
		idStr := strings.TrimPrefix(r.URL.Path, "/reminders/")
		if idStr == "" {
			http.Error(w, "Missing reminder ID", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
			return
		}
		err = store.DeleteReminder(h.DB, id)
		if err == store.ErrNotFound {
			http.Error(w, "Reminder not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to delete reminder", http.StatusInternalServerError)
			return
		} 
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}