#Stage 1
FROM golang:1.23-bullseye AS builder
WORKDIR /smf
COPY .netrc /root/.netrc
COPY . /smf
RUN apt install git wget
RUN go env -w GOPRIVATE=github.com/coranlabs
RUN go mod download

# RUN git clone https://github.com/lakshya-chopra/go /cloudflare-go
# WORKDIR /cloudflare-go/src
# RUN ./make.bash

# ENV GOROOT /cloudflare-go 
# ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
# ENV go /cloudflare-go/bin/go


RUN mkdir -p /go-1.23
RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz \
    && tar -C /go-1.23 -xzf go1.23.0.linux-amd64.tar.gz
ENV GOROOT_BOOTSTRAP=/go-1.23/go


RUN git clone https://github.com/lakshya-chopra/go /cloudflare-go
RUN rm -rf /usr/local/go/src/crypto
RUN rm -rf /usr/local/go/src/vendor
RUN cp /cloudflare-go/src/crypto -R /usr/local/go/src/crypto
RUN cp /cloudflare-go/src/vendor -R /usr/local/go/src/vendor
WORKDIR /usr/local/go/src/

RUN ./make.bash

WORKDIR /smf

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o smf ./main.go

# Stage 2
FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates iputils-ping 
WORKDIR /root/
COPY --from=builder /smf/smf /root/
COPY --from=builder /smf/config/. /root/config/.
COPY .netrc /root/.netrc
RUN chmod +x smf

# Command to run application
CMD ["./smf"]
