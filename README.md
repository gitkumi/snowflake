# Snowflake

Yet another Go REST API application generator.

Features:
- Opinionated with the best practices.
- Simplicity with batteries.
- Idiomatic. Every Gopher loves that word.

## Installation

```sh
go install github.com/gitkumi/snowflake@latest
```

## Quick Start

Here is how to generate an application.

```sh
snowflake new acme
```

Snowflake uses sqlite3 for database.  
If you need a different database, you will need to customize the generated code.  

## Stack

Snowflake is built with these packages. Make sure to check their documentation.

#### Dev

- [air](https://github.com/air-verse/air)

#### Routing

- [Gin](https://gin-gonic.com/)

#### Database

- [sqlc](https://github.com/sqlc-dev/sqlc)
- [goose](https://github.com/pressly/goose)

## Make Commands

- `make dev` - Start the development environment with hot reload.
- `make test` - Run tests.
- `make build` - Build the project.
- `make run` - Run the build.
- `make db` - Check database status.
- `make db.up` - Run database migration.
- `make db.down` - Roll back database migration by 1.
- `make db.create` - Create database.
- `make db.destroy` - Destroy database.
- `make db.reset` - Destroy and create database.
- `make create <table_name>` - Create an empty migration file.
