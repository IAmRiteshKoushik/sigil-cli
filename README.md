# SIGIL: A Certificate Dispatch System

A Go-based system for stamping and dispatching certificates via email using RabbitMQ for queue management and concurrent processing.

## Overview

This system handles three types of certificates for Amrita events:

| Certificate Type | Audience | Nature |
|----------------|---------|--------|
| **Certificate of Participation** | Event participants | External |
| **Certificate of Achievement** | Event winners | External |
| **Certificate of Recognition** | Organizers/volunteers | Internal |

The workflow for each certificate type follows a pipeline pattern:
1. **Process**: Parse CSV and publish student data to a certificate queue
2. **Generate**: Dequeue events, stamp PDFs with names/details, publish to dispatch queue
3. **Send**: Dequeue events and send certificates via email

## Prerequisites

- Go 1.25+
- RabbitMQ 4.1.5+
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) for PDF stamping
- [Just](https://github.com/casey/just) task runner
- Docker (for RabbitMQ)

## Setup

1. **Start RabbitMQ:**
   ```bash
   docker compose up -d
   ```
   Access the management UI at http://localhost:15672 (admin/admin123)

2. **Install dependencies:**
   ```bash
   go mod download
   just
   pip install pdfcpu  # or see pdfcpu installation docs
   ```

3. **Configure the application:**
   Copy `configs/env.toml.sample` to `configs/env.toml` and update:
   ```toml
   rabbitmq_url = "amqp://admin:admin123@localhost:5672/sigil-vhost"
   storage_dir = "storage/"
   smtp_host = "smtp.host.com"
   smtp_port = 587
   smtp_username = "email"
   smtp_password = "app_password"
   ```

4. **Build the application:**
   ```bash
   just build
   # or: go build -o sigil main.go
   ```

## Certificate Types & Workflows

### 1. Certificate of Participation (External)

For event participants.

**Step 1: Process CSV into queue**
```bash
./sigil partprocess storage/data/ESCAPE_ROOM.csv
```
Reads `student_name,student_email` from CSV, extracts event name from filename, and publishes to `cert_{event_name}` queue.

**Step 2: Generate certificates**
```bash
./sigil partcertgen ESCAPE_ROOM templates/participation_cert.pdf
```
- Listens on `cert_ESCAPE_ROOM` queue
- Stamps PDF with student name and event name using `just stamp-participation`
- Saves to `storage/cert/ESCAPE_ROOM/{student_email}.pdf`
- Publishes to `dispatch_internal` queue (for internal processing)

**Step 3: Send certificates via email**
```bash
./sigil partcertsend ESCAPE_ROOM
```
- Listens on `dispatch_ESCAPE_ROOM` queue
- Sends email with certificate attachment to each recipient
- Subject: "Anokha 2026 - Certificate of Participation"

**Participation CSV Format:**
```csv
student_name,student_email
John Doe,john.doe@example.com
```

---

### 2. Certificate of Achievement (External)

For event winners (1st/2nd/3rd place, finalists, etc.).

**Step 1: Process winners CSV**
```bash
./sigil winnerprocess storage/winners.csv
```
Reads `event_name,position,student_name,student_email` and publishes to `cert_winners` queue.

**Step 2: Generate winner certificates**
```bash
./sigil winnergen winners templates/participation_cert.pdf
```
- Listens on `cert_winners` queue
- Stamps PDF with name, event name, and position using `just stamp-winner`
- Saves to `storage/cert/winners/{student_email}.pdf`
- Publishes to `dispatch_winners` queue

**Step 3: Send winner certificates**
```bash
./sigil winnersend winners
```
- Listens on `dispatch_winners` queue
- Sends email with certificate attachment
- Subject: "Anokha 2026 - Certificate of Achievement"

**Winners CSV Format:**
```csv
event_name,position,student_name,student_email
Agamotto Protocol,1st Place,Shanu VS,shanuvs@example.com
AI VERSE HACKATHON,2nd Place,Praveen C S,praveen@example.com
```

---

### 3. Certificate of Recognition (Internal)

For organizers, volunteers, and club heads.

**Step 1: Process internal CSV**
```bash
./sigil internalprocess storage/internal_dept.csv
```
Reads `student_name,student_email,event_name` (3 columns) and publishes to `cert_internal` queue.

**Step 2: Send recognition certificates**
```bash
./sigil internalcertsend <event_name>
```
- Listens on `dispatch_{event_name}` queue
- Reads certificate from `storage/cert/internal/{student_email}.pdf`
- Sends email with certificate attachment
- Subject: "Anokha 2026 - Certificate of Recognition"

**Internal CSV Format:**
```csv
student_name,student_email,event_name
Harshitha Chandrasekar,cb.sc.u4aie23231@cb.students.amrita.edu,Amrita Toastmasters Club – Head
```

---

## Queue Management

### Create queues for events
```bash
./sigil create events.txt
```
Reads event names from file and creates:
- `cert_{event}` queue
- `dispatch_{event}` queue

**events.txt format:**
```
ESCAPE_ROOM
CAMPUS_QUEST
HACKTIDE
```

### Queue naming convention

| Queue Pattern | Purpose |
|-------------|---------|
| `cert_{name}` | Certificate generation input |
| `dispatch_{name}` | Email dispatch input |

---

## Certificate Templates

Place PDF templates in the `templates/` directory:

| Template | Use Case |
|----------|----------|
| `participation_cert.pdf` | Participation & Achievement certificates |
| `recognition_cert.pdf` | Recognition certificates |
| `appreciation_cert.pdf` | (Available for future use) |

### PDF Stamping (via Just)

The system uses `pdfcpu` for stamping via the Just task runner:

```bash
# Stamp participation/recognition certificates
just stamp-participation "<NAME>" "<EVENT NAME>" <input.pdf> <output.pdf>

# Stamp winner certificates with position
just stamp-winner "<NAME>" "<EVENT NAME>" "<POSITION>" <input.pdf> <output.pdf>
```

You can adjust text positioning with offset parameters:
```bash
just stamp-participation "John Doe" "EVENT" input.pdf output.pdf nameoffset="-115" eventnameoffset="50"
```

---

## Directory Structure

```
.
├── configs/
│   ├── env.toml           # Application configuration
│   └── env.toml.sample     # Configuration template
├── storage/
│   ├── cert/              # Generated certificates
│   │   ├── {event_name}/  # Participation certificates
│   │   ├── winners/       # Achievement certificates
│   │   └── internal/      # Recognition certificates
│   └── data/              # Participant CSV files
├── templates/             # PDF certificate templates
├── cmd/                   # CLI command implementations
├── pkg/                   # Shared packages
├── compose.yml            # Docker Compose for RabbitMQ
├── justfile               # Build and stamping tasks
└── sigil                  # Compiled binary
```

---

## All CLI Commands

| Command | Description |
|---------|-------------|
| `sigil create [events-file]` | Create RabbitMQ queues for events |
| `sigil partprocess [csv-file]` | Process participation CSV to cert queue |
| `sigil partcertgen [queue] [template]` | Generate participation certificates |
| `sigil partcertsend [event_name]` | Send participation certificates via email |
| `sigil winnerprocess [csv-file]` | Process winners CSV to cert queue |
| `sigil winnergen [queue] [template]` | Generate winner certificates |
| `sigil winnersend [name]` | Send winner certificates via email |
| `sigil internalprocess [csv-file]` | Process internal CSV to cert queue |
| `sigil internalcertsend [event_name]` | Send recognition certificates via email |

---

## Current Status

- ✅ Certificate of Participation workflow
- ✅ Certificate of Achievement workflow
- ✅ Certificate of Recognition workflow
- ✅ RabbitMQ queue management
- ✅ PDF stamping via pdfcpu
- ✅ Email dispatch via SMTP
