# ============================================================
# Stage 1: Build frontend
# ============================================================
FROM --platform=linux/amd64 node:22-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --ignore-scripts
COPY frontend/ ./
ENV NODE_OPTIONS=--max-old-space-size=4096
ENV VITE_WEB=true
RUN npm run build

# ============================================================
# Stage 2: Build Go backend (web mode)
# ============================================================
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
RUN GOFLAGS="-mod=mod" go mod download

COPY backend/ ./backend/
COPY main_web.go ./

ARG APP_VERSION=1.0.0
RUN CGO_ENABLED=0 GOOS=linux GOFLAGS="-mod=mod" go build -tags web -ldflags "-s -w -X main.version=${APP_VERSION}" -o /app/tinyrdm-server .

# ============================================================
# Stage 3: Runtime (nginx + Go backend)
# ============================================================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata nginx \
    && rm -rf /var/cache/apk/* /tmp/*

# Frontend static files
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# Nginx config
COPY docker/nginx.conf /etc/nginx/http.d/default.conf

# Go backend binary
WORKDIR /app
COPY --from=backend-builder /app/tinyrdm-server .
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8086

ENV PORT=8088
ENV GIN_MODE=release
ENV XDG_CONFIG_HOME=/app

ENTRYPOINT ["/entrypoint.sh"]
