FROM debian:bullseye
#or use alpine


# Copy local Go build (already compiled at /home/ubuntu/go-pq/go-1.24)
COPY ./go-1.24 /usr/local/go

