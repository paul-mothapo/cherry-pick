@echo off
echo ===============================================
echo Database Intelligence with Environment Config
echo ===============================================
echo.
echo This script demonstrates different ways to configure database connections using environment variables.
echo.

echo Option 1: Complete Database URL
echo -------------------------------
echo Example commands:
echo $env:DATABASE_URL = "mysql://root:password@localhost:3306/testdb"
echo $env:DATABASE_URL = "postgres://user:password@localhost:5432/mydb?sslmode=disable"
echo $env:DATABASE_URL = "sqlite3://./test.db"
echo go run examples/env-database-example.go
echo.

echo Option 2: Individual Components
echo --------------------------------
echo Example for MySQL:
echo $env:DB_TYPE = "mysql"
echo $env:DB_HOST = "localhost"
echo $env:DB_PORT = "3306"
echo $env:DB_NAME = "testdb"
echo $env:DB_USER = "root"
echo $env:DB_PASSWORD = "password"
echo go run examples/env-database-example.go
echo.

echo Option 3: MongoDB
echo ------------------
echo $env:MONGODB_URL = "mongodb://localhost:27017/mydb"
echo go run examples/env-database-example.go
echo.

echo To run with your configuration:
echo 1. Set your environment variables using one of the methods above
echo 2. Run: go run examples/env-database-example.go
echo.

echo Running example without configuration to show help...
echo.
go run examples/env-database-example.go
