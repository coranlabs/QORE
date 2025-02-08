## Stage 1: Build Go from source in Alpine
FROM golang:1.24.0-alpine AS builder

RUN apk add --no-cache bash build-base git curl

ENV PATH="/usr/local/go/bin:$PATH"

COPY --from=host_alpine /usr/local/go /usr/local/go-pq

WORKDIR /usr/local/go-pq/src

# Build your Go from source from the alpine builder, will copy this built src to the new stage so that the image size remains minimal.
RUN ./make.bash

## Stage 2:

FROM alpine

RUN apk add --no-cache ca-certificates

# Copy the Go build from the previous stage
COPY --from=builder /usr/local/go-pq /usr/local/go

ENV PATH="/usr/local/go/bin:$PATH"

CMD ["go", "version"]

