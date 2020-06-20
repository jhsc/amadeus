# FROM golang:1.10.3-alpine  AS build-env
# WORKDIR /go/src/github.com/jhsc/amadeus/
# COPY . .
# # Build binary
# RUN CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o goapp ./cmd/main

# FROM alpine:3.7
# WORKDIR /app
# COPY --from=build-env /go/src/github.com/jhsc/amadeus/goapp .
# EXPOSE 8080

# ENTRYPOINT [ "./goapp" ]

FROM golang:1.14-alpine  AS build-env
ENV GO111MODULE=on
WORKDIR /go/src/github.com/jhsc/amadeus

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
EXPOSE 8080
# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o amadeus ./cmd/main.go

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/jhsc/amadeus/amadeus .
EXPOSE 8080

ENTRYPOINT [ "./amadeus" ]

##############################################
#VERSION 2
# FROM golang:1.14-stretch

# ENV GO111MODULE=on
# WORKDIR /go/src/gitlab.com/creativebit/vl

# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# COPY . .
# EXPOSE 8989
# CMD ["go run cmd/api/main.go"]


# VERSION 3
# FROM golang:1.14-stretch AS builder

# ENV GO111MODULE=on
# WORKDIR /src

# COPY . .

# RUN target=/go/pkg/mod,sharing=locked \
#     go test ./... \
#     && CGO_ENABLED=0 go build -a -o /main cmd/api/main.go

# # ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# FROM alpine:3.8 as release

# COPY --from=builder /main /
# COPY --from=builder /src/api /api
# COPY --from=builder /src/migrations /migrations

# EXPOSE 8989
# CMD ["/main"]