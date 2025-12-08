FROM golang:1.25.4-alpine3.22 as build

WORKDIR /app

COPY shared/ ./shared

COPY services/parser ./services/parser

WORKDIR /app/services/parser

RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o parser cmd/main.go

FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates

COPY --from=build /app/services/parser/parser /usr/local/bin/parser

WORKDIR /app

CMD ["parser"]
