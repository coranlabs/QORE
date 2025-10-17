# SPDX-License-Identifier: Apache-2.0
# Copyright 2024 CORAN LABS

FROM golang:1.22 AS builder

WORKDIR /app
RUN apt update \
    && apt install --no-install-recommends -y clang llvm gcc-multilib libbpf-dev \
    && rm -rf /var/lib/apt/lists/*


RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

COPY go.mod go.sum ./
COPY .netrc /root/.netrc
RUN go env -w GOPRIVATE=github.com/coranlabs
RUN go mod download
COPY . /app

ARG BPF_ENABLE_LOG "0"
ARG BPF_ENABLE_ROUTE_CACHE "0"
RUN BPF_CFLAGS="" \
    && if [ "$BPF_ENABLE_LOG" = "1" ]; then BPF_CFLAGS="$BPF_CFLAGS -DENABLE_LOG"; fi \
    && if [ "$BPF_ENABLE_ROUTE_CACHE" = "1" ]; then BPF_CFLAGS="$BPF_CFLAGS -DENABLE_ROUTE_CACHE"; fi \
    && BPF_CFLAGS=$BPF_CFLAGS go generate -v ./ebpf/...
RUN CGO_ENABLED=0 go build -v -o bin/hexa ./main.go

# Build the Go binary
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o hexa ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk update && apk add tcpdump && apk add iptables
WORKDIR /app
#COPY --from=builder /app/hexa .
COPY --from=builder /app/bin/ /app/bin/
COPY --from=builder /app/cmd/docs/swagger.* /app/
#COPY --from=builder /app/cmd/ebpf/zeroentrypoint_bpf.o /app/
# COPY ./entrypoint.sh /app/bin/entrypoint.sh
# COPY ./hexa-upf.yaml /app/config/upfcfg.yaml
RUN chmod +x /app/bin/hexa

EXPOSE 8805
#ENTRYPOINT [ "sh", "./bin/entrypoint.sh" ]
ENTRYPOINT ["sh", "-c", "sleep 3000"]
#CMD ["sh", "-c", "sleep 3000"]

#CMD ["./bin/hexa"]
