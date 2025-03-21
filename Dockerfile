# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

# Build arguments
ARG APP_NAME=app
ARG MAIN_PATH=./cmd/api/main.go
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Install necessary build tools
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Install dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags='-w -s -extldflags "-static"' \
    -tags=netgo,osusergo \
    -o /go/bin/${APP_NAME} \
    ${MAIN_PATH}

# Stage 2: Final stage
FROM scratch

# Build arguments for final stage
ARG APP_NAME=app

# Import from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/${APP_NAME} /${APP_NAME}

# Set environment variables
ENV TZ=UTC \
    APP_USER=appuser \
    APP_UID=1000

# Create a non-root user (even though we're using scratch, it's good practice)
USER 1000

# Command to run
ENTRYPOINT ["/${APP_NAME}"]
