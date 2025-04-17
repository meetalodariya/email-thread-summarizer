# Email Thread Summarizer

Effortlessly streamline customer support email management with Summarizer, from inbox to structured summary.

## Architecture

The application consists of:

- API Server (Go/Echo framework)
- Queue Generator Lambda Function
- PostgreSQL Database
- AWS SQS Queue

## Setup

1. Clone the repository

```bash
git clone https://github.com/meetalodariya/email-thread-summarizer.git
cd email-thread-summarizer
```

2. Install dependencies

```bash
go mod download
```

3. Configure environment variables

```bash
cp .env.example .env
# Edit .env with your configuration
```
