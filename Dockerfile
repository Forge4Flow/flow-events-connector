# Stage 1: Build the backend
FROM golang:1.20 AS backend-builder

WORKDIR /forge4flow
COPY . .
RUN go mod download

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v -o forge4flow-core main.go

# Stage 2: Create the final image
FROM alpine:latest

RUN addgroup -S forge4flow-core && adduser -S forge4flow-core -G forge4flow-core
USER forge4flow-core

WORKDIR /forge4flow
COPY --from=backend-builder /forge4flow/forge4flow-core .

EXPOSE 8000

ENTRYPOINT ["./forge4flow-core"]