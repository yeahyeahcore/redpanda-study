# BUILD
FROM golang:1.20.3-alpine AS build

# Update depences
RUN apk update && apk add --no-cache curl
# Create build directory
RUN mkdir /app/bin -p
RUN mkdir /bin/golang-migrate -p
# Download migrate app
RUN GOLANG_MIGRATE_VERSION=v4.15.1 && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/${GOLANG_MIGRATE_VERSION}/migrate.linux-amd64.tar.gz |\
    tar xvz migrate -C /bin/golang-migrate
# Download health check utility
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.6 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
# Set home directory
WORKDIR /app
# Copy go.mod
ADD go.mod go.sum /app/
# Download go depences
RUN go mod download
# Copy all local files
ADD . /app
# Build app
RUN GOOS=linux go build -o bin ./...



# TEST
FROM alpine:latest AS test

# Install packages
RUN apk --no-cache add ca-certificates
ENV GO111MODULE=off
# Create home directory
WORKDIR /app
# Copy build file
COPY --from=build /app/bin/app ./app
# CMD
CMD [ "./app" ]


# MIGRATION
FROM alpine:latest AS migration

# Install packages
RUN apk --no-cache add ca-certificates
# Create home directory
WORKDIR /app
# Copy migration dir
COPY --from=build /app/migrations/test ./migrations
# Install migrate tool
COPY --from=build /bin/golang-migrate /usr/local/bin



# PRODUCTION
FROM alpine:latest AS production

# Install packages
RUN apk --no-cache add ca-certificates
# Create home directory
WORKDIR /app
# Copy build file
COPY --from=build /app/bin/app ./app
# Copy grpc health probe dir
COPY --from=build /bin/grpc_health_probe /bin/grpc_health_probe
# Install migrate tool
COPY --from=build /bin/golang-migrate /usr/local/bin
# CMD
CMD ["./app"]
