# SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.22.5-bullseye AS builder

LABEL maintainer="ONF <omec-dev@opennetworking.org>"

#RUN apt remove cmdtest yarn
RUN apt-get update
RUN apt-get -y install apt-transport-https ca-certificates
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg > pubkey.gpg
RUN apt-key add pubkey.gpg
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" |  tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update
RUN apt-get -y install gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev  nodejs yarn libssl-dev git

# WORKDIR /cloudflare-go/src
# COPY ~/cloudflare-go /cloudflare-go/

# COPY ./cloudflare-go /cloudflare-go

#due to git clone issues
RUN git config --global http.postBuffer 524288000
RUN git config --global core.compression 9 repack

RUN git clone https://github.com/lakshya-chopra/go /cloudflare-go
WORKDIR /cloudflare-go/src
RUN ./make.bash
RUN apt-get clean

# or use GOPATH/src
# RUN echo "alias go=/cloudflare-go/bin/go" >> ~/.bashrc
#RUN $go version
#ENV GOFLAGS=''
#RUN go env
# ENV GOINSECURE *

ENV GOROOT /cloudflare-go
ENV PATH $GOROOT/bin:$PATH
RUN echo $GOPATH

RUN echo "alias go=/cloudflare-go/bin/go" >> ~/.bashrc

ENV GOFLAGS=-mod=vendor



RUN go env > /tmp/go_env_output.txt && cat /tmp/go_env_output.txt
ENV GOSUMDB off
# RUN go env -w GOSUMDB=off

ENV CGO_CFLAGS="-I/usr/include/openssl"
ENV CGO_LDFLAGS="-L/usr/lib/x86_64-linux-gnu -lssl -lcrypto"

#RUN apt-get install -y ca-certificates && update-ca-certificates
#RUN go env -w GOPROXY=direct
#ENV GONOSUMDB=github.com/omec-project/*

#RUN go clean -modcache
# RUN go get github.com/omec-project/nrf@v1.4.1


RUN cd $GOPATH/src && mkdir -p amf
COPY . $GOPATH/src/amf
RUN cd $GOPATH/src/amf
# RUN go mod tidy

# RUN go get github.com/lakshya-chopra/http2_util
RUN cd $GOPATH/src/amf && make all


# FROM alpine:3.8 as amf

# # RUN mkdir -p /usr/lib64
# # RUN mkdir -p /usr/include/openssl

# # COPY --from=builder /usr/lib/x86_64-linux-gnu/* /usr/lib64
# # COPY --from=builder /usr/include/openssl/ /usr/include/openssl

# RUN apk add --no-cache openssl-dev musl-dev



# ENV CGO_CFLAGS="-I/usr/include/openssl"
# ENV CGO_LDFLAGS="-L/usr/lib64/ -lssl -lcrypto"

# ENV ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.22

# LABEL description="ONF open source 5G Core Network" \
#     version="Stage 3"

# ARG DEBUG_TOOLS

# # Install debug tools ~ 100MB (if DEBUG_TOOLS is set to true)
# RUN apk update
# RUN apk add -U vim strace net-tools curl netcat-openbsd bind-tools bash

# # Set working dir
# WORKDIR /free5gc
# RUN mkdir -p amf/

# # Copy executable and default certs
# COPY --from=builder /cloudflare-go/src/amf/bin/* ./amf
# COPY --from=builder /cloudflare-go/src/amf/cert/* /free5gc/support/TLS/
# WORKDIR /free5gc/amf


#Deploy stage using Ubuntu instead of Alpine
FROM ubuntu:20.04 as amf

RUN apt-get update && apt-get install -y \
    libssl-dev dnsutils bash


ENV CGO_CFLAGS="-I/usr/include/openssl"
ENV CGO_LDFLAGS="-L/usr/lib/x86_64-linux-gnu -lssl -lcrypto"
ENV ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.22

LABEL description="ONF open source 5G Core Network" \
    version="Stage 3"

ARG DEBUG_TOOLS

# Set working directory and copy the AMF binary and certificates from the builder stage
WORKDIR /free5gc
RUN mkdir -p amf/

COPY --from=builder /go/src/amf/bin/* ./amf
COPY --from=builder /go/src/amf/cert/* /free5gc/support/TLS/

WORKDIR /free5gc/amf
