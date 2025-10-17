#Stage 1
FROM golang:1.23-alpine AS builder
WORKDIR /coran_console
COPY .netrc /root/.netrc
COPY . /coran_console/
RUN apk add --no-cache git
RUN go env -w GOPRIVATE=github.com/coranlabs
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o coran_console ./main.go

# Stage 2
FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates iputils-ping 
WORKDIR /root/
COPY --from=builder /coran_console/coran_console /root/
COPY --from=builder /coran_console/config/CORAN_CONSOLE.yaml . /root/config/
COPY .netrc /root/.netrc
RUN chmod +x coran_console

# Command to run application
CMD ["./coran_console"]

# WebUI uses the port 5000
EXPOSE 5000
# BillingDomain uses the port 2121 as FTP
EXPOSE 2121