# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy the whole module. go.sum is generated during the build because the
# sandbox that authored this code had no network access to pre-resolve it.
COPY . .

ENV GOFLAGS=-mod=mod
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

# ---- Runtime stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/server /app/server
COPY --from=builder /out/seed /app/seed

# Runtime data (uploads, backups, exports) lives here; mount a volume in prod.
RUN mkdir -p /app/data && chown -R app:app /app
USER app

ENV RENTACAR_ENV=production \
    RENTACAR_HOST=0.0.0.0 \
    RENTACAR_PORT=5000 \
    RENTACAR_DATA_DIR=/app/data

EXPOSE 5000

ENTRYPOINT ["/app/server"]
