FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Generate templ first: the *_templ.go files must exist before Tailwind
# runs, because tailwind/tailwind.config.js scans ./internal/view/**/*.go
# to pick up class strings emitted by templ.
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate

# Build CSS. The -c flag is required: the config lives in tailwind/, not the
# repo root, so without it Tailwind falls back to an empty default config,
# scans no content, and purges every utility class (yielding a ~5KB stub).
RUN apk add --no-cache curl && \
    curl -sL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 -o /usr/local/bin/tailwindcss && \
    chmod +x /usr/local/bin/tailwindcss && \
    tailwindcss -c ./tailwind/tailwind.config.js -i tailwind/input.css -o web/static/css/site.css --minify

RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates

COPY --from=build /bin/server /bin/server
COPY --from=build /src/web /web

# SQLite database location. The app creates the parent dir if missing,
# but we make it explicit so an orchestrator (compose, k8s, etc.) knows
# which path to mount a persistent volume at.
ENV DB_PATH=/var/lib/app/app.db
RUN mkdir -p /var/lib/app
VOLUME ["/var/lib/app"]

EXPOSE 8080
ENTRYPOINT ["/bin/server"]
