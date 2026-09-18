# syntax=docker/dockerfile:1

# ---------- Tailwind CSS ----------
FROM node:22-alpine AS css
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@10.28.2 --activate
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/css ./web/css
COPY internal/view ./internal/view
RUN mkdir -p web/static && pnpm css

# ---------- Go build ----------
FROM golang:1.27-alpine AS build
WORKDIR /app
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN templ generate
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gonshin ./cmd/gonshin

# ---------- Runtime ----------
FROM alpine:3.21
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /gonshin /usr/local/bin/gonshin
COPY --from=build /app/web/static ./web/static
COPY --from=css /app/web/static/app.css ./web/static/app.css

ENV \
    ADDR=":8510" \
    HOYOLAB_API_BASE="https://hoyo.dotcchix.dev" \
    CACHE_TTL_SECONDS="1200" \
    PUBLIC_BASE_URL=""

EXPOSE 8510
USER app
ENTRYPOINT ["gonshin"]
