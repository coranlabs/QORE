# Define build and target base images
ARG BUILD_IMAGE=golang:1.23.0-alpine
ARG TARGET_IMAGE=alpine:3.18

# Build stage: use Go to compile the binaries for all components
FROM $BUILD_IMAGE AS builder
WORKDIR /nwdaf
ENV CGO_ENABLED=0

# Copy go modules and download dependencies
COPY go.* ./
RUN go mod download

# Copy the rest of the source code
COPY . ./

# Build binaries for each application component
RUN go build -o bin/coranlabs_sbi ./Application_entity/cmd/sbi/main.go
RUN go build -o bin/coranlabs_events ./Application_entity/cmd/events/main.go
RUN go build -o bin/coranlabs_analytics ./Application_entity/cmd/analytics/main.go
RUN go build -o bin/coranlabs_engine ./Application_entity/cmd/engine/main.go
RUN go build -o bin/coranlabs_nbiml ./Application_entity/cmd/nbiml/main.go

# Final stage: create a minimal production image
FROM $TARGET_IMAGE

# Update and install necessary tools
RUN apk update && apk add -U net-tools curl bash

# Copy the built binaries and .env file from the builder stage
COPY --from=builder /nwdaf/bin/coranlabs_sbi /nwdaf/coranlabs_sbi
COPY --from=builder /nwdaf/bin/coranlabs_events /nwdaf/coranlabs_events
COPY --from=builder /nwdaf/bin/coranlabs_analytics /nwdaf/coranlabs_analytics
COPY --from=builder /nwdaf/bin/coranlabs_engine /nwdaf/coranlabs_engine
COPY --from=builder /nwdaf/bin/coranlabs_nbiml /nwdaf/coranlabs_nbiml
COPY config/nwdaf.env /nwdaf/config/nwdaf.env

# Set the working directory
WORKDIR /nwdaf

# Expose the correct port based on the service
ARG SERVICE_NAME
# Choose the port depending on the service
# Ensure that each service is associated with its proper port
EXPOSE 8080  
EXPOSE 8081  
EXPOSE 8082  
EXPOSE 8083  
EXPOSE 8084

# Add a shell script to select the binary based on the SERVICE_NAME environment variable
COPY entrypoint.sh /nwdaf/entrypoint.sh
RUN chmod +x /nwdaf/entrypoint.sh

# Use the entrypoint script to run the appropriate binary
CMD ["/nwdaf/entrypoint.sh"]
