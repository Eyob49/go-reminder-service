package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	handlers "reminder/internal/handler"
	"reminder/internal/scheduler"
	"reminder/store"

	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


func main(){
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API key loaded: %v", os.Getenv("RESEND_API_KEY") != "")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
	log.Fatal("Failed to connect to database:", err)
    }

    defer db.Close()

    err = store.InitDB(db)
    if err != nil {
	     log.Fatal("Failed to initialize database:", err)
    }
    handler := &handlers.ReminderHandler{DB: db}
    ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduler.StartScheduler(ctx, db)
	http.HandleFunc("/reminders", handler.HandleReminders)
	http.HandleFunc("/reminders/", handler.HandleReminders)
    log.Printf("Server is running on http://localhost:%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}