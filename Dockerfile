FROM golang:1.22.3 as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=1 go build -o shounic-best-clips .

FROM debian:12.8-slim
WORKDIR /app
COPY --from=build /app/shounic-best-clips .
COPY --from=build /app/www .

CMD ["./shounic-best-clips"]
