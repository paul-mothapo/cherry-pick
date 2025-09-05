# Enterprise Database Intelligence System

**Go 1.21+** | A comprehensive database analysis and intelligence platform for MySQL, PostgreSQL, SQLite, and MongoDB.

## Overview

The Enterprise Database Intelligence System is a modular Go application that provides comprehensive database analysis, monitoring, security assessment, and optimization capabilities. Built with advanced Go architecture patterns including dependency injection and interface-based design, it offers actionable insights to help enterprises optimize their database infrastructure.

## Key Features

- **🔍 Comprehensive Analysis**: Tables, columns, indexes, constraints, and relationships
- **🛡️ Security Assessment**: PII detection, vulnerability scanning, access control analysis
- **📊 Performance Monitoring**: Real-time monitoring, alerting, and optimization suggestions
- **🔗 Data Lineage**: Dependency tracking and impact analysis
- **⚡ Query Optimization**: SQL analysis with performance recommendations
- **📈 Automated Reporting**: Scheduled analysis with customizable alerts

## Supported Databases

- MySQL 5.7+
- PostgreSQL 10+
- SQLite 3.x
- MongoDB 4.0+

## Quick Start

### Web UI (New!)

**Development:**
```bash
# Terminal 1 - Start UI development server
cd web
npm install
npm run dev

# Terminal 2 - Start Go server
go mod tidy
go run cmd/server/main.go
```

The web interface will be available at `http://localhost:3000`

### Command Line

```go
service, err := intelligence.CreateSimpleService("mysql", "user:pass@tcp(localhost:3306)/db")
if err != nil {
    log.Fatal(err)
}
defer service.Close()

report, err := service.AnalyzeDatabase()
```

## Documentation

- **[Running the Application](RUNNING.md)** - Setup, configuration, and execution guide
- **[Contributing](CONTRIBUTING.md)** - Development and contribution guidelines
- **[Security Policy](SECURITY.md)** - Security reporting and policies

## License

MIT License - see [LICENSE](LICENSE) file for details.