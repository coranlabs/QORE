FROM golang:1.24.0 AS builder

# Install build dependencies and bootstrap Go
#RUN apt update && apt install -y bash build-essential git curl

ENV PATH="/usr/local/go/bin:$PATH"

COPY --from=host_bullseye /usr/local/go /usr/local/go-pq

WORKDIR /usr/local/go-pq/src

# Build your Go from source from the alpine builder, will copy this built src to the new stage so that the image size remains minimal.
RUN ./make.bash

## Stage 2:

FROM debian:bullseye

# Copy the Go build from the previous stage
COPY --from=builder /usr/local/go-pq /usr/local/go

ENV PATH="/usr/local/go/bin:$PATH"

CMD ["go", "version"]

