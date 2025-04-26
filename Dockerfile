FROM golang:1.21 as builder
WORKDIR /usr/src/app
COPY . /usr/src/app/

# Set build arguments for version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Build with version information
RUN go build -o jetspotter \
    -ldflags "-linkmode external -extldflags -static \
              -X jetspotter/internal/version.Version=${VERSION} \
              -X jetspotter/internal/version.Commit=${COMMIT} \
              -X jetspotter/internal/version.BuildTime=${BUILD_TIME}" \
    cmd/jetspotter/jetspotter.go


FROM alpine:latest
COPY --from=builder /usr/src/app/jetspotter .
CMD [ "./jetspotter" ]
