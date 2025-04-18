# SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.21-bullseye AS builder

LABEL maintainer="ONF <omec-dev@opennetworking.org>"

#RUN apt remove cmdtest yarn
RUN apt-get update && apt-get -y install apt-transport-https ca-certificates
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg > pubkey.gpg
RUN apt-key add pubkey.gpg
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" |  tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get -y install gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev  nodejs yarn
RUN apt-get clean


RUN cd $GOPATH/src && mkdir -p pcf
COPY . $GOPATH/src/pcf
RUN cd $GOPATH/src/pcf \
    && make all

FROM alpine:3.16 as pcf

LABEL description="ONF open source 5G Core Network" \
    version="Stage 3"

ARG DEBUG_TOOLS

# Install debug tools ~ 100MB (if DEBUG_TOOLS is set to true)
RUN apk update && apk add -U vim strace net-tools curl netcat-openbsd bind-tools

# Set working dir
WORKDIR /free5gc
RUN mkdir -p pcf/

# Copy executable and default certs
COPY --from=builder /go/src/pcf/bin/* ./pcf
COPY --from=builder /go/src/pcf/cert/* /free5gc/support/TLS/
WORKDIR /free5gc/pcf
