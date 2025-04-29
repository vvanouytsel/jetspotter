FROM golang:1.23 AS builder
WORKDIR /usr/src/app
COPY . /usr/src/app/

# Set build arguments for version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Build with CGO disabled for a fully static binary
ENV CGO_ENABLED=0

# Build with version information
RUN go build -o jetspotter -ldflags "-X jetspotter/internal/version.Version=${VERSION} -X jetspotter/internal/version.Commit=${COMMIT} -X jetspotter/internal/version.BuildTime=${BUILD_TIME}" cmd/jetspotter/jetspotter.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/jetspotter .
CMD ["./jetspotter"]
