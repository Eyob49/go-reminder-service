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
	godotenv.Load("../.env")
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
    log.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}