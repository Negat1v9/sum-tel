FROM golang:1.25.3-alpine3.22 as build

WORKDIR /app

COPY shared ./shared

COPY services/migrator ./services/migrator

WORKDIR /app/services/migrator

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

RUN go build -o migrator ./cmd/main.go

FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates

COPY --from=build /app/services/migrator/migrator /usr/local/bin/migrator

WORKDIR /services/migrator

CMD ["migrator"]
