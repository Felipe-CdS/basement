# check=error=true

# =============================================================================
# BASE: Common development dependencies
# =============================================================================

ARG GO_IMG_TAG=1.25-alpine3.23

FROM golang:${GO_IMG_TAG} AS dev-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# =============================================================================
# PROD-BUILD
# =============================================================================

FROM dev-base AS prod-build

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
# DEV-HOT
# =============================================================================
FROM dev-base AS dev-hot

RUN go install github.com/air-verse/air@v1.61.1

WORKDIR /app

CMD ["air", "-c", ".air.toml"]
