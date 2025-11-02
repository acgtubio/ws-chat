FROM golang:1.23-alpine3.19 AS base
WORKDIR /chat/

# System dependencies
RUN apk update && \
    apk add --no-cache ca-certificates tzdata git && \
    update-ca-certificates

# Install air and delve
RUN go install github.com/cosmtrek/air@v1.49.0

### Run the backend project
FROM base AS backend
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENTRYPOINT ["air"]
