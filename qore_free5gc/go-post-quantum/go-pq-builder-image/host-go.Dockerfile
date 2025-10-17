FROM alpine:3.19

# Copy your local Go build (already compiled at /home/ubuntu/go-pq/go-1.24)
COPY ./go-1.24 /usr/local/go

