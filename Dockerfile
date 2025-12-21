# check=error=true

# =============================================================================
# GENERATE STAGE

FROM ghcr.io/a-h/templ:latest AS generate-stage

COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

# =============================================================================
# FETCH STAGE

FROM golang:1.25-alpine3.23 AS fetch-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# =============================================================================
# PROD-BUILD

FROM fetch-stage AS prod-build

ARG BUILD_GOOS
ENV BUILD_GOOS=${BUILD_GOOS}

COPY . .

RUN GOOS=${BUILD_GOOS} GOARCH=arm64 go build -mod=mod -ldflags "-w -s" -o ./tmp/basement ./cmd/

FROM alpine:3.23 AS dev-static-runtime
LABEL org.opencontainers.image.source=https://nugu.dev/basement

WORKDIR /app

RUN apk --no-cache add ca-certificates
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

COPY --from=prod-build /app/tmp/basement .
COPY --from=prod-build /app/config.yaml .

ENTRYPOINT ["/app/basement"]

# =============================================================================
# HOT-BUILD

FROM fetch-stage AS hot-build

RUN go install github.com/air-verse/air@v1.61.1

WORKDIR /app

CMD ["air", "-c", ".air.toml"]

# =============================================================================
# TAILWIND-WATCH
FROM debian:bookworm-slim AS tailwind-watch

ENV TAILWIND_URL="https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.18/tailwindcss-linux-arm64"

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl watchman && rm -rf /var/lib/apt/lists/*

RUN curl -L -o /usr/local/bin/tailwindcss ${TAILWIND_URL}
RUN chmod +x /usr/local/bin/tailwindcss
