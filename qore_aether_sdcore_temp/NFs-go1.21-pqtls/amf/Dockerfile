# SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.21-bullseye AS builder

LABEL maintainer="ONF <omec-dev@opennetworking.org>"

#RUN apt remove cmdtest yarn
RUN apt-get update
RUN apt-get -y install apt-transport-https ca-certificates
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg > pubkey.gpg
RUN apt-key add pubkey.gpg
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" |  tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update
RUN apt-get -y install gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev  nodejs yarn git

RUN git clone https://github.com/lakshya-chopra/go /cloudflare-go
WORKDIR /cloudflare-go/src
RUN ./make.bash
RUN apt-get clean

# ENV GOROOT /cloudflare-go
ENV GOPATH /cloudflare-go
ENV PATH $GOPATH/bin:$PATH
    
    # or use GOPATH/src
ENV go /cloudflare-go/bin/go
ENV GOFLAGS=''

RUN cd $GOPATH/src && mkdir -p amf
COPY . $GOPATH/src/amf
RUN cd $GOPATH/src/amf \
    && make all

RUN echo $PWD

FROM alpine:3.8 as amf

ENV ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.22

LABEL description="ONF open source 5G Core Network" \
    version="Stage 3"

ARG DEBUG_TOOLS

# Install debug tools ~ 100MB (if DEBUG_TOOLS is set to true)
RUN apk update
RUN apk add -U vim strace net-tools curl netcat-openbsd bind-tools bash

# Set working dir
WORKDIR /free5gc
RUN mkdir -p amf/

# Copy executable and default certs
COPY --from=builder /cloudflare-go/src/amf/bin/* ./amf
COPY --from=builder /cloudflare-go/src/amf/cert/* /free5gc/support/TLS/
WORKDIR /free5gc/amf
