# SIGIL: A Certificate Dispatch System

A Go-based system for processing CSV data and sending personalized 
certificates via email using RabbitMQ for queue management and 
concurrent processing.

## Overview

This system reads student data from CSV files, validates eligibility, 
creates event-specific RabbitMQ queues, and processes certificates through 
a concurrent workflow. Each event gets dedicated `cert_` and `dispatch_` 
queues for certificate processing and email delivery.

## Setup

1. **Start RabbitMQ:**
   ```bash
   docker-compose up -d
   ```

2. **Build the application:**
   ```bash
   go build -o sigil .
   ```

3. **Create queues for events:**
   ```bash
   ./sigil create events.txt
   ```

## Usage

The CLI tool supports the following commands:

### Queue Management
```bash
./sigil create [events-file]
```
Reads events from the specified file and creates:
- `cert_{event}` queues for certificate processing
- `dispatch_{event}` queues for email dispatch

### CSV Processing
```bash
./sigil process [csv-file]
```
Processes a single CSV file and publishes student data to the appropriate event queue.

### Batch Processing
```bash
./sigil process-batch [reports-folder]
./sigil process-batch reports --move
```
Processes all CSV files in the specified folder:
- Scans folder recursively for `.csv` files
- Shows progress with file-by-file status
- Optional `--move` flag moves processed files to `processed/` subfolder
- Provides detailed success/failure summary

## Directory Structure

```
.
├── reports/          # Place CSV files here for batch processing
│   ├── processed/   # Automatically created when using --move flag
│   └── .gitkeep
├── templates/       # Certificate template files
│   └── .gitkeep
├── config.toml      # RabbitMQ and application configuration
└── sigil           # Compiled binary
```

## CSV File Format

CSV files should have the following columns:
- `student_name` - Full name of the student
- `student_email` - Email address for certificate delivery
- `event_name` - Event name (optional, defaults to filename)

Example:
```csv
student_name,student_email
John Doe,john.doe@example.com
Jane Smith,jane.smith@example.com
```

## Current Status

Phase 1-2 implementation complete with:
- ✅ RabbitMQ container setup
- ✅ CLI queue management
- ✅ Event file processing
- ✅ Dynamic queue creation per event
- ✅ CSV parsing and validation
- ✅ Single file processing
- ✅ Batch processing with progress tracking

Phases 3-4 (Certificate generation, email dispatch) in development.
