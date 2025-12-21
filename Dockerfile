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
