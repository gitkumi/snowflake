# Snowflake

Yet another Go web application generator.

Features:
- Opinionated with the best practices.
- Simplicity with batteries.
- Idiomatic. Every Gopher loves that word.

## Installation

```sh
go install github.com/gitkumi/snowflake@latest
```

## Quick Start

Here is how to generate a web application with sqlite3 as database.

```sh
snowflake new -t api -d sqlite3 acme
```

## Flags

- **type (-t)**
  - `web` - Generates a web application that serves HTML. (Requires `Node.js 20.x` and `pnpm`)
  - `api` - Generates a RESTful API application.

- **database (-d)**
  - `sqlite3`
  - `postgres`
  - `mysql`

## Stack

Snowflake is built with these packages. Make sure to check their documentation.

#### Dev

- [air](https://github.com/air-verse/air)

#### Routing

- [Gin](https://gin-gonic.com/)

#### Database

- [sqlc](https://github.com/sqlc-dev/sqlc)
- [goose](https://github.com/pressly/goose)

#### Front-end (Web only)

- [templ](https://templ.guide/)
- [Tailwind](https://tailwindcss.com/)
- [esbuild](https://esbuild.github.io/)
- [tygo](https://github.com/gzuidhof/tygo)

## Make Commands

- `make dev` - Start the development environment with hot reload.
- `make test` - Run tests.
- `make build` - Build the project.
- `make run` - Run the build.
- `make db` - Check database status.
- `make up` - Run database migration.
- `make down` - Roll back database migration by 1.
- `make reset` - Roll back all migration.
- `make create <table_name>` - Create an empty migration file.
