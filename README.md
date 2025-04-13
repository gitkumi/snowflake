# Snowflake

Yet another Go application generator.

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
snowflake new acme -d postgres
```

## Flags

- `-d`: Database type (`sqlite3`, `postgres`, or `mysql`)
- `-t`: App type (`api` or `web`)

Check out `help` for more flags.

## Stack

Snowflake is built with these packages. Make sure to check their documentation.

#### Dev

- [air](https://github.com/air-verse/air)

#### Routing

- [Gin](https://gin-gonic.com/)

#### Database

- [sqlc](https://github.com/sqlc-dev/sqlc)
- [goose](https://github.com/pressly/goose)

#### Test Runner

- [gotestum](https://github.com/gotestyourself/gotestsum)

#### Templating

- [templ](https://templ.guide) (only for "web" app type)
