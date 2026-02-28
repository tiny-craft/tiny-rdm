# ============================================================
# Stage 1: Build frontend
# ============================================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
ENV NODE_OPTIONS=--max-old-space-size=4096
ENV VITE_WEB=true
RUN npm run build

# ============================================================
# Stage 2: Build Go backend (web mode)
# ============================================================
FROM golang:1.24-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GOFLAGS=-mod=mod

COPY backend/ ./backend/
COPY main.go main_web.go ./
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -tags web -ldflags "-s -w -X main.version=1.2.6" -o /app/tinyrdm .

# ============================================================
# Stage 3: Runtime
# ============================================================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata fontconfig font-noto font-noto-cjk

WORKDIR /app
COPY --from=backend-builder /app/tinyrdm .
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8088

ENV PORT=8088
ENV GIN_MODE=release
ENV XDG_CONFIG_HOME=/app

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./tinyrdm"]
