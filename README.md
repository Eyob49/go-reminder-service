# go-reminder-service

A REST API backend in Go that schedules and sends reminders via email, SMS, or both.

## Features

- Create reminders with a title, message, send time, email, phone, and delivery channel
- Store reminders in SQLite
- Periodically check for due reminders and dispatch them automatically
- Support channels: `email`, `sms`, `both`

## Tech stack

- Go
- SQLite
- Resend (email delivery)
- Twilio (SMS delivery)

## Project structure

- `/cmd/main.go` — application entrypoint
- `/internal/handler/reminder.go` — HTTP request handling
- `/internal/model/model.go` — reminder model definition
- `/internal/notifier/email.go` — email delivery logic
- `/internal/notifier/sms.go` — SMS delivery logic
- `/internal/scheduler/scheduler.go` — periodic reminder dispatch
- `/store/reminder_store.go` — SQLite persistence

## Setup

### Prerequisites

- Go installed
- Twilio account and credentials
- Resend account and API key

### Installation

1. Clone the repository:

```bash
git clone https://github.com/Eyob49/go-reminder-service.git
cd go-reminder-service
```

2. Install dependencies:

```bash
go mod tidy
```

3. Create a `.env` file in the project root and add your credentials.

4. Run the server:

```bash
go run ./cmd
```

The server starts on `http://localhost:8080` and creates `reminders.db` in the current directory.

## Environment variables

Create a `.env` file in the repository root with the following values:

```env
RESEND_API_KEY=
TWILIO_ACCOUNT_SID=
TWILIO_AUTH_TOKEN=
TWILIO_PHONE_NUMBER=
```

## API Endpoints

### POST /reminders

Create a new reminder.

Request body example:

```json
{
  "title": "Doctor appointment",
  "message": "Don't forget your 3pm appointment",
  "email": "your@email.com",
  "phone": "+2519********",
  "channel": "both",
  "sendAt": "2026-05-27T12:00:00Z"
}
```

Important notes:

- `sendAt` must be a future UTC datetime in RFC3339 format, for example: `2026-05-27T12:00:00Z`
- `channel` must be one of `email`, `sms`, or `both`
- If `channel` is `email`, `email` is required
- If `channel` is `sms`, `phone` is required
- If `channel` is `both`, both `email` and `phone` are required

Response:

- `201 Created` with the created reminder payload

### GET /reminders

Retrieve all reminders.

Response:

- `200 OK` with a JSON array of reminders

### DELETE /reminders/{id}

Delete a reminder by ID.

Example:

```bash
delete http://localhost:8080/reminders/1
```

Response:

- `204 No Content`

## Notes

- The scheduler checks for due reminders once per minute and marks them as sent after dispatch.
- The SQLite database is initialized automatically if it does not exist.
