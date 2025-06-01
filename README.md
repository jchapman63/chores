# Chore Rotation App

A simple command-line application to manage chore rotations for roommates, built with Go.

## Overview

This application helps roommates track and rotate chores on a weekly basis. The system stores each roommate's information along with their assigned chore and automatically rotates chores each week following the pattern: BATHROOM → FLOOR → COUNTER.

## Features

- Simple roommate database with name, phone number, and chore
- Automatic weekly chore rotation
- SMS notifications via AWS SNS
- Command-line interface for managing roommates and chores
- Designed to be run via cron job for weekly rotations

## Project Structure

```
chores/
├── cmd/                  # Application entrypoints
│   └── chores/           # Main CLI application
├── internal/             # Internal packages
│   ├── db/               # Database code
│   │   ├── schema/       # SQL schema definitions
│   │   ├── query/        # SQL queries for SQLc
│   │   └── sqlc/         # Generated Go code (via SQLc)
│   ├── rotation/         # Chore rotation logic
│   └── notification/     # SMS notification service
├── config/               # Configuration files
├── go.mod                # Go module file
├── sqlc.yaml             # SQLc configuration
└── README.md             # This file
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- PostgreSQL database
- AWS account (for SMS notifications)

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
4. Initialize the database:
   ```
   go run cmd/chores/main.go init
   ```
5. Configure AWS credentials for SMS notifications

### Command-Line Usage

The application supports the following commands:

```
Usage:
  chores init                       - Initialize the database schema
  chores list                       - List all roommates and their chores
  chores add <name> <phone> <chore> - Add a new roommate
  chores rotate                     - Rotate chores
  chores notify                     - Send chore notifications
  chores weekly                     - Perform weekly rotation and notification
```

### Setting up the Weekly Rotation

To set up the weekly rotation job, add the following to your crontab:

```
# Run the chore rotation every Monday at 8:00 AM
0 8 * * 1 cd /path/to/chores && go run cmd/chores/main.go weekly >> /var/log/chores_rotation.log 2>&1
```

## Chore Rotation Logic

Chores follow this rotation sequence:
1. BATHROOM
2. FLOOR
3. COUNTER

Each week, roommates move to the next chore in the sequence.