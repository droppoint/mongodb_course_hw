
FROM golang:1.25 AS builder
WORKDIR /app
COPY loadgenerator/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o loadgenerator main.go
    
FROM scratch
WORKDIR /app
COPY --from=builder /app/loadgenerator /app/loadgenerator

ENTRYPOINT ["/app/loadgenerator"]