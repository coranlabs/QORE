# Stage 1
FROM golang:1.22.6-bullseye AS builder
WORKDIR /nssf
COPY .netrc /root/.netrc
COPY . /nssf
RUN apt install git
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

WORKDIR /nssf

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o nssf ./main.go

# Stage 2
FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates iputils-ping 
WORKDIR /root/
COPY --from=builder /nssf/nssf  ./
COPY --from=builder /nssf/config/CORAN_NSSF.yaml ./config/CORAN_NSSF.yaml
COPY .netrc /root/.netrc
RUN chmod +x nssf

# Command to run application
CMD ["./nssf"]
