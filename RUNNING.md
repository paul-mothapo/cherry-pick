# Running the Application

This guide provides comprehensive instructions for setting up and running the Enterprise Database Intelligence System.

## Prerequisites

- **Go 1.21+** (required)
- Access to a database for testing (optional - you can use the included SQLite example)

## Installation

1. **Navigate to the project directory:**
   ```bash
   cd cherry-pick
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Verify Go version:**
   ```bash
   go version  # Should show Go 1.21 or higher
   ```

## Configuration

### Method 1: Environment Variables (Recommended)

Set database connection using environment variables:

**PowerShell (Windows):**
```powershell
# MySQL
$env:DATABASE_URL = "mysql://root:password@localhost:3306/mydb"

# PostgreSQL
$env:DATABASE_URL = "postgres://user:password@localhost:5432/mydb?sslmode=disable"

# SQLite
$env:DATABASE_URL = "./database.db"

# MongoDB
$env:MONGODB_URL = "mongodb://localhost:27017/mydb"
```

**Bash (Linux/macOS):**
```bash
# MySQL
export DATABASE_URL="mysql://root:password@localhost:3306/mydb"

# PostgreSQL
export DATABASE_URL="postgres://user:password@localhost:5432/mydb?sslmode=disable"

# SQLite
export DATABASE_URL="./database.db"

# MongoDB
export MONGODB_URL="mongodb://localhost:27017/mydb"
```

### Method 2: Individual Components

**PowerShell:**
```powershell
$env:DB_TYPE = "mysql"
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"
$env:DB_NAME = "mydb"
$env:DB_USER = "root"
$env:DB_PASSWORD = "password"
```

**Bash:**
```bash
export DB_TYPE="mysql"
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_NAME="mydb"
export DB_USER="root"
export DB_PASSWORD="password"
```

### Method 3: Configuration File (Optional)

You can create a `config.json` file for advanced configuration:
```json
{
  "database_connections": {
    "primary": "mysql://user:password@localhost:3306/mydb"
  },
  "analysis_settings": {
    "sample_size": 1000
  }
}
```

## Running the Application

### Quick Demo

```bash
# Run the main demonstration
go run main.go
```

### Database Analysis Examples

#### 1. Environment-Based Analysis
```bash
# Set your database URL
export DATABASE_URL="mysql://user:pass@localhost:3306/mydb"

# Run analysis
go run examples/env-database-example.go
```

#### 2. Direct Database Connection
```bash
go run examples/real-database-example.go
```

#### 3. MongoDB Analysis
```bash
export MONGODB_URL="mongodb://localhost:27017/mydb"
go run examples/mongodb-example.go
```

#### 4. Create Test Database
```bash
# Create a test SQLite database with sample data
go run examples/create-test-db.go

# Then analyze it
export DATABASE_URL="./test.db"
go run examples/env-database-example.go
```

### Using the Batch Script (Windows)

```batch
# Run the environment configuration helper
run-with-env.bat
```

## Usage Examples

### Basic Analysis

```go
package main

import (
    "fmt"
    "log"
    "github.com/cherry-pick/pkg/intelligence"
)

func main() {
    // Create service
    service, err := intelligence.CreateSimpleService(
        "mysql", 
        "user:password@tcp(localhost:3306)/dbname",
    )
    if err != nil {
        log.Fatal(err)
    }
    defer service.Close()

    // Analyze database
    report, err := service.AnalyzeDatabase()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Database: %s\n", report.DatabaseName)
    fmt.Printf("Tables analyzed: %d\n", len(report.Tables))
    fmt.Printf("Health Score: %.2f\n", report.Summary.HealthScore)
}
```

### Security Analysis

```go
// Check for security issues
issues, err := service.AnalyzeSecurity()
if err != nil {
    log.Fatal(err)
}

for _, issue := range issues {
    fmt.Printf("SECURITY: %s - %s\n", issue.Title, issue.Description)
}
```

### Query Optimization

```go
// Optimize a query
suggestion, err := service.OptimizeQuery("SELECT * FROM users WHERE active = 1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Original: %s\n", suggestion.OriginalQuery)
fmt.Printf("Optimized: %s\n", suggestion.OptimizedQuery)
```

### Monitoring Setup

```go
import "time"

// Schedule periodic analysis
service.ScheduleAnalysis(24*time.Hour, func(report *types.DatabaseReport) {
    fmt.Printf("Scheduled analysis completed: Health Score %.2f\n", report.Summary.HealthScore)
    
    // Check for alerts
    alerts, _ := service.CheckAlerts()
    for _, alert := range alerts {
        if alert.Severity == "high" {
            fmt.Printf("ALERT: %s\n", alert.Message)
        }
    }
})
```

## Database Setup (Optional)

If you want to test with your own database, you'll need appropriate access. The application only requires SELECT permissions for basic analysis.

## Output Formats

### JSON Export
```go
jsonData, err := service.ExportReport(report, "json")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(jsonData))
```

### Save to File
```go
import "os"

jsonData, _ := service.ExportReport(report, "json")
err := os.WriteFile("report.json", jsonData, 0644)
if err != nil {
    log.Fatal(err)
}
```

## Troubleshooting

### Common Issues

#### "No database configuration found"
- Ensure environment variables are set correctly
- Check that `DATABASE_URL` or `MONGODB_URL` is defined
- Verify the connection string format

#### Connection Errors
- Check your database connection string format
- Verify database server is running and accessible
- Ensure credentials are correct

#### Permission Errors
- Ensure database user has SELECT permissions
- For SQLite, check file permissions and directory access

#### Go Module Issues
```bash
# Re-download dependencies if needed
go mod download
```

### Debug Mode

Enable verbose logging:
```go
import "log"

// Set log level for debugging
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### Performance Tips

For large databases, you can adjust sample size in configuration if needed.

## Environment Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | Complete SQL database connection string | `mysql://user:pass@host:port/db` |
| `MONGODB_URL` | Complete MongoDB connection string | `mongodb://localhost:27017/mydb` |
| `DB_TYPE` | Database type | `mysql`, `postgres`, `sqlite3` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `3306`, `5432` |
| `DB_NAME` | Database name | `myapp` |
| `DB_USER` | Database username | `root`, `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_SSLMODE` | SSL mode (PostgreSQL) | `disable`, `require` |
| `TEST_DB_NAME` | Test database filename | `test.db` |

## Next Steps

- Check [Contributing Guidelines](CONTRIBUTING.md) for development information  
- Review [Security Policy](SECURITY.md) for security considerations
