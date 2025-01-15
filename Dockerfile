FROM golang:1.22.3 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=1 go build \
    -ldflags "-X main.commitSHA=$(git rev-parse HEAD)" \
    -o shounic-best-clips \
    .

FROM debian:12.8-slim
WORKDIR /app
COPY --from=build /app/shounic-best-clips .

CMD ["./shounic-best-clips"]
