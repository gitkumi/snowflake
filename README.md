# Snowflake

Yet another Go REST API application generator.

Features:
- Opinionated with the best practices.
- Simplicity with batteries.
- Idiomatic. Every Gopher loves that word.

## Prerequisite

snowflake requires git and the following Go libraries.

```sh
go install github.com/air-verse/air@latest &&
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && 
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Installation


```sh
go install github.com/gitkumi/snowflake@latest
```

## Quick Start

Here is how to generate an application.

```sh
snowflake new acme -d postgres
```

## Flags

- `-d`: Database type (`sqlite3`, `postgres`, or `mysql`)

## Stack

Snowflake is built with these packages. Make sure to check their documentation.

#### Dev

- [air](https://github.com/air-verse/air)

#### Routing

- [Gin](https://gin-gonic.com/)

#### Database

- [sqlc](https://github.com/sqlc-dev/sqlc)
- [goose](https://github.com/pressly/goose)

## Commands

- `make dev` - Start the development environment with hot reload.
- `make test` - Run tests.
- `make build` - Build the project.
- `make run` - Run the build.
- `make audit` - Audit the project.
- `make tidy` - Tidy the modules and format the project.
- `make db` - Check database status.
- `make db.up` - Run database migration.
- `make db.down` - Roll back database migration by 1.
- `make db.create` - Create database.
- `make db.destroy` - Destroy database.
- `make db.reset` - Destroy and create database.
- `make create <table_name>` - Create an empty migration file.

## TODO

- dev dashboards (telemetry, admin, email dev, file storage dev)
- consider migrating dev commands in go instead of makefile (e.g. db migrate stuff)
- db: consider supporting gorm/std lib
- router: consider adding choice/std db
- serve html choice
