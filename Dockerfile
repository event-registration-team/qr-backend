FROM golang:1.26-alpine AS build

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/event-registration ./cmd

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=build /out/event-registration /app/event-registration
COPY migrations /app/migrations

USER app
EXPOSE 8080

CMD ["/app/event-registration"]
