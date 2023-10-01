FROM golang:1.21 as builder
WORKDIR /usr/src/app
COPY . /usr/src/app/
RUN go build -o jetspotter -ldflags "-linkmode external -extldflags -static" cmd/jetspotter/jetspotter.go


FROM alpine:latest
COPY --from=builder /usr/src/app/jetspotter .
CMD [ "./jetspotter" ]
