# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.26-bookworm AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY api/ ./api/
COPY db/ ./db/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /cmd/pin-pyre-db-api ./cmd/pin-pyre-db-api

# Stage 2: Run the tests
FROM build AS test
RUN go test ./...

# Stage 3: Create the release image
FROM debian:bookworm-slim AS release

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system app \
    && useradd --system --gid app --uid 10001 app

WORKDIR /app

RUN mkdir -p /app/data && chown -R app:app /app
COPY --from=build --chown=app:app /cmd/pin-pyre-db-api /cmd/pin-pyre-db-api
COPY --chown=app:app migrations/ ./migrations/

EXPOSE 8080

ENV PORT=8080
ENV DB_PATH=/app/data/tyres.db

VOLUME ["/app/data"]

USER app

ENTRYPOINT ["/cmd/pin-pyre-db-api"]
