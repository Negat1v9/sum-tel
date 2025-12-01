FROM golang:1.25.4-alpine3.22 as build

WORKDIR /app

COPY shared/ ./shared

COPY services/core ./services/core

WORKDIR /app/services/core

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

RUN go build -o api cmd/main.go

FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates

COPY --from=build /app/services/core/api /usr/local/bin/api

WORKDIR /app

CMD ["api"]
