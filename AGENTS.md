# Certificate Dispatch System

## Overview
Build a Go-based certificate dispatch system that processes CSV data, 
sends personalized certificates via email using RabbitMQ for queue management 
and concurrent processing.

## Requirements
1. **Certificate Template Integration**: Media team provides certificate template file that needs parsing and data injection
2. **CSV Data Processing**: Parse CSV with columns: event_name, student_name, student_email
3. **Eligibility Filtering**: Only process students with at least one attended schedule
4. **Dynamic Queue Management**: Create separate RabbitMQ queues per event based on events.txt, removing scrapped events
5. **Events Configuration**: Use events.txt file containing all event names to initialize queues if they don't exist
6. **Folder Structure**: Create dedicated folders for each event by event name to store generated certificates
7. **Concurrent Processing**: Spawn goroutines equal to active event count (e.g., 50 events = 50 goroutines)
8. **Certificate Workflow**: Extract payload → Parse certificate template → Inject student data → Generate certificate → Save to event folder → Send email → Acknowledge completion

## Implementation Plan

### Phase 1: Core Infrastructure
- Set up Go project structure with proper modules
- Configure RabbitMQ connection and queue management
- Implement CSV parser with validation logic
- Create certificate template parser and data injection system
- Set up SMTP email client configuration
- Implement folder management for event-specific certificate storage

### Phase 2: Queue System
- Parse events.txt to get list of all valid events
- Implement dynamic queue creation per event based on events.txt
- Initialize queues if they don't exist on startup
- Build message payload structure (event_id, student_data, template_vars)
- Create queue consumer/producer patterns
- Add error handling and retry mechanisms
- Implement queue cleanup for completed/failed messages

### Phase 3: Concurrent Processing
- Design goroutine pool management system
- Implement graceful shutdown handling
- Add rate limiting for email sending
- Create monitoring and logging for processing status
- Build acknowledgment system for message completion

### Phase 4: Data Flow
1. **CSV Ingestion**: Parse and validate input data
2. **Eligibility Check**: Filter students with attended schedules
3. **Queue Distribution**: Route students to event-specific queues
4. **Template Processing**: Parse certificate template and inject student data
5. **Certificate Generation**: Create personalized certificate files
6. **File Storage**: Save generated certificates to event-specific folders
7. **Email Dispatch**: Send certificates via SMTP with attachments
8. **Acknowledgment**: Mark messages as processed and cleanup

### Technical Specifications
- **Language**: Go 1.19+
- **Queue**: RabbitMQ with amqp library
- **Email**: SMTP client (net/smtp or external service)
- **Template**: Certificate template parser with variable injection
- **CSV**: encoding/csv for data parsing
- **Concurrency**: Goroutines with worker pool pattern
- **Config**: Environment-based configuration management

### Error Handling
- CSV validation and malformed data handling
- SMTP connection failures and retry logic
- RabbitMQ connection recovery
- Template rendering error handling
- Graceful degradation for partial failures

### Monitoring & Logging
- Processing statistics (success/failure rates)
- Queue depth monitoring
- Email delivery status tracking
- Performance metrics (processing time, throughput)
- Error categorization and alerting

### Security Considerations
- Email credential management (environment variables)
- Input sanitization for template injection
- Rate limiting to prevent spamming
- Access control for queue operations
- Data privacy compliance for student information

### Deployment Architecture
- Docker containerization
- Kubernetes deployment with horizontal scaling
- Environment-specific configurations
- Database migrations (if needed for tracking)
- Backup and disaster recovery procedures
