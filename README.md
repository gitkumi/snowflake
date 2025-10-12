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

Here is how to create an application.

```sh
# tui
snowflake new

# cli
snowflake run acme -d postgres
```

## Stack

Snowflake is built with these packages. Make sure to check their documentation.

#### Dev

- [air](https://github.com/air-verse/air)

#### Routing

- [Gin](https://gin-gonic.com/)

#### Database

Supports "sqlite3", "postgres", "mysql"

- [sqlc](https://github.com/sqlc-dev/sqlc)
- [goose](https://github.com/pressly/goose)

#### Templating 

- [templ](https://templ.guide) 


#### Testing

- [gotestsum](https://github.com/gotestyourself/gotestsum)
