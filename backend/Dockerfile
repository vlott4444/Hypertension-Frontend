FROM golang:1.27-alpine AS builder
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /out/migrate ./cmd/migrate && \
    CGO_ENABLED=0 go build -o /out/seed ./cmd/seed && \
    CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM alpine:3.22
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/ /app/
COPY backend/resources/styles /app/resources/styles
COPY frontend/public/media /app/resources/media
COPY backend/templates /app/templates
EXPOSE 8080
CMD ["/bin/sh","-c","/app/migrate && /app/seed && exec /app/server"]
