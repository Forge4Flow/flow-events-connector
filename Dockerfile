# Stage 1: Build the backend
FROM golang:1.20 AS backend-builder

WORKDIR /flow-events-connector
COPY . .
RUN go mod download

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v -o flow-events-connector main.go

# Stage 2: Create the final image
FROM alpine:latest

RUN addgroup -S flow-events-connector && adduser -S flow-events-connector -G flow-events-connector
USER flow-events-connector

WORKDIR /forge4flow
COPY --from=backend-builder /forge4flow/flow-events-connector .

ENTRYPOINT ["./flow-events-connector"]