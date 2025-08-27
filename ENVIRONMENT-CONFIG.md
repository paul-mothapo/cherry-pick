# Environment Variable Configuration

This document explains how to configure database connections using environment variables instead of hardcoding connection strings.

## Quick Start

### Option 1: Complete Database URLs

Set a single environment variable with the complete connection string:

```powershell
# MySQL
$env:DATABASE_URL = "mysql://root:password@localhost:3306/testdb"

# PostgreSQL
$env:DATABASE_URL = "postgres://user:password@localhost:5432/mydb?sslmode=disable"

# SQLite
$env:DATABASE_URL = "./database.db"

# MongoDB
$env:MONGODB_URL = "mongodb://localhost:27017/mydb"
$env:MONGODB_URL = "mongodb+srv://user:password@cluster.mongodb.net/mydb"
```

### Option 2: Individual Components

Set individual environment variables for each database component:

```powershell
# MySQL Components
$env:DB_TYPE = "mysql"
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"
$env:DB_NAME = "testdb"
$env:DB_USER = "root"
$env:DB_PASSWORD = "password"

# PostgreSQL Components
$env:DB_TYPE = "postgres"
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_NAME = "mydb"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "password"
$env:DB_SSLMODE = "disable"

# SQLite Components
$env:DB_TYPE = "sqlite3"
$env:DB_NAME = "./database.db"
```

## Running Examples

After setting your environment variables, run:

```powershell
# New environment-based example (recommended)
go run examples/env-database-example.go

# Original example (still works)
go run examples/real-database-example.go
```

## Environment Variables Reference

### Complete URLs

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | Complete SQL database connection string | `mysql://user:pass@host:port/db` |
| `MONGODB_URL` | Complete MongoDB connection string | `mongodb://localhost:27017/mydb` |

### Individual Components

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `DB_TYPE` | Database type | - | `mysql`, `postgres`, `sqlite3` |
| `DB_HOST` | Database host | `localhost` | `localhost`, `db.example.com` |
| `DB_PORT` | Database port | varies by type | `3306`, `5432`, `27017` |
| `DB_NAME` | Database/schema name | - | `myapp`, `production` |
| `DB_USER` | Database username | - | `root`, `postgres`, `admin` |
| `DB_PASSWORD` | Database password | - | `password123` |
| `DB_SSLMODE` | SSL mode (PostgreSQL) | `disable` | `disable`, `require`, `verify-full` |

### Test Database

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `TEST_DB_NAME` | Test SQLite database filename | `test.db` | `mytest.db` |

## Examples by Database Type

### MySQL

```powershell
# Complete URL
$env:DATABASE_URL = "mysql://root:password@localhost:3306/myapp"

# Components
$env:DB_TYPE = "mysql"
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"
$env:DB_NAME = "myapp"
$env:DB_USER = "root"
$env:DB_PASSWORD = "password"

go run examples/env-database-example.go
```

### PostgreSQL

```powershell
# Complete URL
$env:DATABASE_URL = "postgres://postgres:password@localhost:5432/myapp?sslmode=disable"

# Components
$env:DB_TYPE = "postgres"
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_NAME = "myapp"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "password"
$env:DB_SSLMODE = "disable"

go run examples/env-database-example.go
```

### SQLite

```powershell
# Complete URL
$env:DATABASE_URL = "./myapp.db"

# Components
$env:DB_TYPE = "sqlite3"
$env:DB_NAME = "./myapp.db"

go run examples/env-database-example.go
```

### MongoDB

```powershell
# Local MongoDB
$env:MONGODB_URL = "mongodb://localhost:27017/myapp"

# MongoDB Atlas
$env:MONGODB_URL = "mongodb+srv://username:password@cluster.mongodb.net/myapp"

# With authentication
$env:MONGODB_URL = "mongodb://user:password@localhost:27017/myapp"

go run examples/env-database-example.go
```

## Creating Test Data

Create a test SQLite database with sample data:

```powershell
# Default test database
go run examples/create-test-db.go

# Custom database name
$env:TEST_DB_NAME = "mytest.db"
go run examples/create-test-db.go

# Then analyze it
$env:DATABASE_URL = "./mytest.db"
go run examples/env-database-example.go
```

## Configuration Priority

The system checks for configuration in this order:

1. `DATABASE_URL` environment variable (highest priority)
2. `MONGODB_URL` environment variable
3. Individual component variables (`DB_TYPE`, `DB_HOST`, etc.)
4. If none found, shows configuration help

## Security Best Practices

1. **Never commit environment files** - Add `.env` files to `.gitignore`
2. **Use least privilege** - Create database users with minimal required permissions
3. **Use SSL/TLS** - Enable SSL for production databases
4. **Rotate passwords** - Regularly change database passwords
5. **Use secrets management** - In production, use proper secrets management systems

## Troubleshooting

### No Configuration Found

If you see "No database configuration found", ensure you've set at least one of:
- `DATABASE_URL` or `MONGODB_URL`
- `DB_TYPE` + `DB_NAME` (minimum for individual components)

### Connection Errors

- **MySQL**: Check if MySQL server is running and credentials are correct
- **PostgreSQL**: Verify PostgreSQL service is running and SSL settings
- **SQLite**: Ensure the database file path is correct and accessible
- **MongoDB**: Confirm MongoDB service is running and connection string format

### Permission Errors

- Ensure the database user has necessary permissions (SELECT, INSERT, etc.)
- For SQLite, check file system permissions on the database file

## Advanced Usage

### Multiple Databases

You can analyze different databases by changing environment variables:

```powershell
# Analyze production database
$env:DATABASE_URL = "postgres://user:pass@prod-db:5432/app"
go run examples/env-database-example.go

# Switch to staging
$env:DATABASE_URL = "postgres://user:pass@staging-db:5432/app"
go run examples/env-database-example.go
```

### Batch Scripts

Create batch scripts for different environments:

```batch
@echo off
echo Setting up Production Database Analysis
set DATABASE_URL=mysql://user:pass@prod-server:3306/production
go run examples/env-database-example.go
```

This makes it easy to switch between different database environments without changing code.
