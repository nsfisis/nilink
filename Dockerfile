FROM golang:1.25 AS builder

WORKDIR /app

RUN apt-get update && \
    apt-get install -y libsqlite3-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o nilink .

##########################################
FROM gcr.io/distroless/cc-debian12

WORKDIR /app
COPY --from=builder /app/nilink /app

EXPOSE 8080
ENTRYPOINT ["/app/nilink"]
CMD ["serve", "-addr", ":8080"]
