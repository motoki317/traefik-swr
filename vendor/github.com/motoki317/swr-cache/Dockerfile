FROM --platform=$BUILDPLATFORM golang:1-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED 0

COPY ./go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/server -ldflags "-s -w" ./cmd/server

FROM alpine:3 as runner

WORKDIR /app

COPY --from=builder /app/server ./
ENTRYPOINT ["/app/server"]
