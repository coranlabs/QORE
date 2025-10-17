# Stage 1
FROM golang:1.23.0-alpine AS builder
WORKDIR /scp
COPY .netrc /root/.netrc
COPY . /scp
RUN apk add --no-cache git
RUN go env -w GOPRIVATE=github.com/coranlabs
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o scp ./main.go

# Stage 2
FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates iputils-ping net-tools 
WORKDIR /root/
# COPY --from=builder /scp/scp /scp/config/CORAN_scp.yaml ./
COPY --from=builder /scp/scp ./
COPY --from=builder /scp/config/CORAN_SCP.yaml ./config/
# COPY .netrc /root/.netrc
RUN chmod +x scp

# Command to run application
CMD ["./scp"]
