# STEP 1: Build an executable binary
FROM golang:alpine AS builder

# Install Git to fetch dependencies
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

# Fetch dependencies using go mod (go 1.11)
RUN go mod download
RUN go mod verify

# Build binary
RUN go build -o /go/bin/reliefexchange_backend

# STEP 2: Build image solely with executable
FROM scratch
COPY --from=builder /go/bin/reliefexchange_backend /go/bin/reliefexchange_backend

# Import necessary data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Set up environment variables
ENV GIN_MODE release
ENV PORT 5050

# Expose port
EXPOSE 5050

# Run backend
ENTRYPOINT ["/go/bin/reliefexchange_backend"]