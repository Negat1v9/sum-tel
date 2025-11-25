FROM golang:1.25.4-alpine3.22 as build

WORKDIR /app

COPY shared/ ./shared

COPY services/narrator ./services/narrator

WORKDIR /app/services/narrator

RUN go mod download

RUN go build cmd/main.go

FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates

COPY --from=build /app/services/narrator/main /usr/local/bin/main

WORKDIR /app

CMD ["main"]
