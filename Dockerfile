FROM golang:1.22-bookworm AS build

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential pkg-config libasound2-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o my-theme ./...

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends libasound2 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/my-theme /app/my-theme
COPY --from=build /app/static /app/static
COPY --from=build /app/theme.mp3 /app/theme.mp3

EXPOSE 8080

ENTRYPOINT ["/app/my-theme"]

