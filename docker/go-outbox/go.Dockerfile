FROM golang:1.24 AS base

FROM base AS builder

WORKDIR /build
COPY . ./
RUN go build -mod=mod ./cmd/outbox

FROM base AS final

ARG PORT

WORKDIR /app
COPY --from=builder /build/outbox /build/.env ./
COPY --from=builder /build/outbox ./

CMD ["/app/outbox"]