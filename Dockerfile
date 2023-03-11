#syntax=docker/dockerfile:1.2

# Prepare phase
FROM --platform=$BUILDPLATFORM golang:1.15.15-alpine3.14 as base

RUN apk --update --no-cache add build-base

# Build phase
FROM base as build

ARG FLAGS TARGETOS TARGETARCH
ARG SERVICE_NAME=ddns-service

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY "./cmd/${SERVICE_NAME}/main.go" /app
COPY internal /app/internal

RUN ls -la /app

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -v -ldflags="${FLAGS}" -trimpath \
    -o /out/cloudflare-ddns main.go

# Finalize phase
FROM alpine:3.14

ENV CF_API_KEY key
ENV CF_API_EMAIL email
ENV DOMAINS domains
ENV AUTH_USER user
ENV AUTH_PASS pass

VOLUME /app
WORKDIR /app

COPY --from=build /out/cloudflare-ddns /usr/bin/

EXPOSE 8008

ENTRYPOINT ["sh", "-c", "/usr/bin/cloudflare-ddns --cf-api-key=${CF_API_KEY} --cf-api-email=${CF_API_EMAIL} --domains=${DOMAINS} --auth-user=${AUTH_USER} --auth-password=${AUTH_PASS}"]