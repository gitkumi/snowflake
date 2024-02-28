# Snowflake

Opinionated Go web application generator.  

## Quick Start

Here is how to generate a web project using sqlite3 as database.

```sh
snowflake new -t web -d sqlite3 acme
```

### Flags

- **type (-t)**
  - `web`: Generates a web application that serves HTML.
  - `api`: Generates RESTful API application.

- **database (-d)**
  - `sqlite3`
  - `postgres`
  - `mysql`
