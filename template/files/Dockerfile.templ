FROM golang:1.21-bookworm AS builder

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/gzuidhof/tygo@latest
RUN apt-get update && apt-get install -y make curl gnupg
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - &&\
    apt-get install -y nodejs
RUN npm install -g pnpm

WORKDIR /app

COPY . .

RUN pnpm install
RUN make build

FROM golang:1.21-alpine AS server

ARG PORT=8080
ENV PORT=${PORT}

ARG DATABASE_URL
ENV DATABASE_URL=${DATABASE_URL}
ENV MIGRATE_DB=true

ENV ENVIRONMENT=production
ENV GIN_MODE=release

WORKDIR /app

COPY --from=builder /app/bin/main /app/main

EXPOSE 8080

CMD ["./main"]
