FROM golang:1.24-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/server/

FROM node:22-alpine AS frontend-builder

WORKDIR /app

COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ .
RUN npm run build

FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    font-noto \
    font-noto-cjk

WORKDIR /app

COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/dist ./web/dist

RUN mkdir -p /app/data /app/uploads /app/generated

ENV PORT=8080
ENV DATABASE_PATH=/app/data/snapreport.db
ENV UPLOAD_DIR=/app/uploads
ENV PDF_DIR=/app/generated
ENV ENVIRONMENT=production

EXPOSE 8080

CMD ["./server"]
