# Stage 1: Build the Go amflication
FROM golang:1.23.0-bullseye AS builder

# Set the working directory inside the container
WORKDIR /amf

# Copy the current directory contents into the container at /amf
COPY . /amf
COPY .netrc /root/.netrc
RUN ls -aL
RUN apt install git wget 
RUN apt update

RUN go env -w GOPRIVATE=github.com/coranlabs
# Download Go module dependencies
RUN go mod download

COPY go1.23.0.linux-amd64.tar.gz /tmp/

RUN mkdir -p /go-1.23 && \
    tar -C /go-1.23 -xzf /tmp/go1.23.0.linux-amd64.tar.gz && \
    rm /tmp/go1.23.0.linux-amd64.tar.gz

ENV GOROOT_BOOTSTRAP=/go-1.23/go


RUN git clone https://github.com/lakshya-chopra/go /cloudflare-go
#COPY cloudflare-go /cloudflare-go
RUN rm -rf /usr/local/go/src/crypto
RUN rm -rf /usr/local/go/src/vendor
RUN cp /cloudflare-go/src/crypto -R /usr/local/go/src/crypto
RUN cp /cloudflare-go/src/vendor -R /usr/local/go/src/vendor

# RUN rm -rf /usr/local/go/src/
# RUN cp /cloudflare-go/src/ -R /usr/local/go/src/
WORKDIR /usr/local/go/src/

# ENV GOROOT /usr/local/go
RUN ./make.bash

# ENV GOROOT /cloudflare-go 
# ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
# ENV go /cloudflare-go/bin/go

# RUN /usr/local/go/bin/go
RUN go version

WORKDIR /amf

# Build the Go amflication for Linux with CGO disabled
# Adjust the path if your main.go is in the root directory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o amf ./main.go

# Stage 2: Create a lightweight image for running the amflication
FROM alpine:latest

# Install CA certificates to handle HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the built amflication from the builder stage
COPY --from=builder /amf/amf ./
COPY --from=builder /amf/config/CORAN_AMF.yaml ./config/

# Make the amflication executable
RUN chmod +x amf
RUN ls -al

# Command to run your amflication
CMD ["./amf"]
