FROM golang:1.10.3-alpine  AS build-env
WORKDIR /go/src/github.com/jhsc/amadeus/
COPY . .
# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o goapp ./cmd/main

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/jhsc/amadeus/goapp .
EXPOSE 4020

ENTRYPOINT [ "./goapp" ]