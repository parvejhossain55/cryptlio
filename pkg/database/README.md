# Cryplio Database Package

Shared database utilities, connection pool, migration runner, and schema definitions.

## Directory Structure

```
pkg/database/
├── migrations/          # SQL migration files (see /migrations/README.md)
│   ├── 000_types.sql
│   ├── 001_core_tables.sql
│   ├── ...
│   └── 016_seed_data.sql
├── connection.go        # Database connection pool setup
├── models/              # Go structs for database entities
│   ├── user.go
│   ├── trade.go
│   └── ...
├── repository/          # Data access layer (CRUD operations)
│   ├── user_repository.go
│   ├── trade_repository.go
│   └── ...
├── query/               # Prepared SQL queries (if needed)
└── README.md

</directories>