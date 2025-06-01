# Chore Rotation App

A simple application to manage chore rotations for roommates, built with Go.

## Overview

This application helps roommates track and rotate chores on a weekly basis. The system stores each roommate's information along with their assigned chore and automatically rotates chores every Monday at 9am following the pattern: BATHROOM → FLOOR → COUNTER.

## Features

- Simple roommate database with name, phone number, and chore
- Automatic weekly chore rotation every Monday at 9am
- Built-in scheduler using cron expressions
- PostgreSQL database for data persistence
- Graceful shutdown handling

## Project Structure

```
chores/
├── cmd/                  # Application entrypoints
│   └── chores/           # Main application
├── internal/             # Internal packages
│   ├── db/               # Database code
│   │   ├── schema/       # SQL schema definitions
│   │   ├── query/        # SQL queries for SQLc
│   │   └── sqlc/         # Generated Go code (via SQLc)
│   └── rotation/         # Chore rotation logic
├── config/               # Configuration files
├── go.mod                # Go module file
├── sqlc.yaml             # SQLc configuration
└── README.md             # This file
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL database

### Setup

1. Clone the repository
2. Install dependencies:
   ```
   go mod tidy
   ```
3. Generate database code with SQLc:
   ```
   sqlc generate
   ```
4. Set up your database schema (using PostgreSQL):
   ```sql
   CREATE TABLE roommates (
     id SERIAL PRIMARY KEY,
     name TEXT NOT NULL,
     phone TEXT NOT NULL,
     chore TEXT NOT NULL
   );
   ```
5. Configure database connection by setting environment variables (or use defaults):
   ```
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=chores
   ```

### Running the Application

Run the application:

```
go run cmd/chores/main.go
```

The application will:
1. Connect to the PostgreSQL database using the configured connection details
2. Initialize the rotation service with database queries
3. Schedule a chore rotation job for every Monday at 9am
4. Wait for a termination signal (Ctrl+C) to shut down gracefully

The application uses an embedded cron scheduler, so there's no need to set up external cron jobs. As long as the application is running, it will handle the chore rotations automatically.

### Deployment

For production deployment, build the application:

```
go build -o chores ./cmd/chores
```

You can then run the binary directly:

```
./chores
```

Consider using a process manager like systemd to ensure the application stays running.

#### Example systemd Service File

Create a file at `/etc/systemd/system/chores.service`:

```
[Unit]
Description=Chores Rotation Application
After=network.target postgresql.service

[Service]
Type=simple
User=chores
WorkingDirectory=/opt/chores
ExecStart=/opt/chores/chores
Restart=on-failure
Environment=DB_HOST=localhost
Environment=DB_PORT=5432
Environment=DB_USER=postgres
Environment=DB_PASSWORD=postgres
Environment=DB_NAME=chores

[Install]
WantedBy=multi-user.target
```

Then enable and start the service:

```
sudo systemctl enable chores
sudo systemctl start chores
```

## Chore Rotation Logic

Chores follow this rotation sequence:
1. BATHROOM
2. FLOOR
3. COUNTER

Each week, roommates move to the next chore in the sequence. The rotation happens automatically every Monday at 9am when the application is running.

## Database Schema

The application uses a simple database schema with a single table:

```sql
CREATE TABLE roommates (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  phone TEXT NOT NULL,
  chore TEXT NOT NULL
);
```

Where:
- `id`: Unique identifier for each roommate
- `name`: Roommate's name
- `phone`: Roommate's phone number
- `chore`: Current assigned chore (one of: "BATHROOM", "FLOOR", "COUNTER")

## Logging and Troubleshooting

The application outputs logs to standard output (stdout) and standard error (stderr), including:

- Application startup information
- Database connection status
- Scheduler initialization
- Chore rotation execution and results
- Application shutdown events

### Common Issues

1. **Database Connection Issues**:
   - Verify PostgreSQL is running: `pg_isready -h <host> -p <port>`
   - Check credentials and permissions
   - Ensure the database exists: `psql -h <host> -U <user> -c "SELECT 1 FROM pg_database WHERE datname = 'chores'"`

2. **Missing Tables**:
   - Confirm the roommates table exists: `psql -h <host> -U <user> -d chores -c "\dt"`
   - Create the table if missing using the SQL schema provided above

3. **Scheduler Not Running**:
   - The application must remain running for the scheduler to work
   - Check system logs if using systemd: `journalctl -u chores.service`

For more detailed troubleshooting, you can run the application with verbose logging:

```
DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=chores ./chores 2>&1 | tee chores.log
```