# check=error=true

# =============================================================================
# GENERATE TEMPL STAGE

FROM ghcr.io/a-h/templ:v0.3.1020 AS templ-generate-stage

COPY --chown=65532:65532 . /app
WORKDIR /app

RUN ["templ", "generate"]

# =============================================================================
# GENERATE TAILWIND STAGE

FROM --platform=$BUILDPLATFORM coutito/tailwindcss:v4.3.3 AS tailwind-generate-stage 

WORKDIR /app
COPY --from=templ-generate-stage /app .

RUN tailwindcss -i ./assets/css/tailwind_input.css -o ./assets/css/tailwind_styles.css --minify

# =============================================================================
# PROD-BUILD

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build-stage

WORKDIR /app
COPY --from=tailwind-generate-stage /app .

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=mod -ldflags "-w -s" -o ./tmp/basement ./cmd/

# =============================================================================
# FINAL-IMAGE

FROM alpine:3.23 AS final-image

LABEL org.opencontainers.image.source=https://nugu.dev/basement

WORKDIR /app

RUN mkdir /assets
COPY --from=tailwind-generate-stage /app/assets ./assets/
COPY --from=build-stage /app/tmp/basement .
COPY --from=build-stage /app/tokens .

# TURSO NEEDS: libgcc gcompat
RUN apk --no-cache add ca-certificates libgcc gcompat

ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

EXPOSE 3000
ENTRYPOINT ["/app/basement"]
