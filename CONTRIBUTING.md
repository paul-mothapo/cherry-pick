# Contributing to Enterprise Database Intelligence System

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Documentation](#documentation)

## Code of Conduct

This project follows a standard code of conduct:

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect different viewpoints and experiences
- Show empathy towards other community members

## Getting Started

### Prerequisites

- **Go 1.21+** (required)
- Git
- A supported database (MySQL, PostgreSQL, SQLite, or MongoDB) for testing
- Familiarity with Go modules and interfaces

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/cherry-pick.git
   cd cherry-pick
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/original-org/cherry-pick.git
   ```

## Development Setup

### Environment Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Set up test database:**
   ```bash
   # Create a test SQLite database
   go run examples/create-test-db.go
   
   # Set environment variable
   export DATABASE_URL="./test.db"
   ```

3. **Verify setup:**
   ```bash
   go run main.go
   go run examples/env-database-example.go
   ```

### IDE Configuration

**VS Code:**
- Install the Go extension
- Configure `settings.json`:
  ```json
  {
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint"
  }
  ```

**GoLand/IntelliJ:**
- Enable Go modules support
- Configure code style to use `gofmt`

## Project Structure

```
cherry-pick/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── LICENSE                 # License file
├── README.md               # Main documentation
├── RUNNING.md              # Running instructions
├── CONTRIBUTING.md         # This file
├── SECURITY.md             # Security policy
├── run-with-env.bat        # Windows batch script
├── pkg/                    # Core packages
│   ├── types/             # Data structures
│   ├── interfaces/        # Interface definitions
│   ├── intelligence/      # Main service layer
│   ├── analyzer/          # Database analysis
│   ├── connector/         # Database connections
│   ├── insights/          # Intelligence generation
│   ├── monitoring/        # Monitoring and alerts
│   ├── security/          # Security analysis
│   ├── optimization/      # Query optimization
│   ├── config/           # Configuration management
│   └── utils/            # Utility functions
└── examples/              # Example applications
    ├── create-test-db.go
    ├── env-database-example.go
    ├── example.go
    ├── mongodb-example.go
    └── real-database-example.go
```

### Package Responsibilities

- **`types/`**: Core data structures and models
- **`interfaces/`**: Interface definitions for dependency injection
- **`intelligence/`**: Main service coordination and factory patterns
- **`analyzer/`**: Database schema and data analysis
- **`connector/`**: Database connection management
- **`insights/`**: Intelligence generation and reporting
- **`monitoring/`**: Alerting, scheduling, and comparison
- **`security/`**: Security analysis and PII detection
- **`optimization/`**: Query analysis and optimization
- **`config/`**: Configuration management
- **`utils/`**: Shared utility functions

## Coding Standards

### Go Style Guidelines

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Naming Conventions

- **Packages**: Short, lowercase, single word when possible
- **Functions**: CamelCase, start with uppercase for exported functions
- **Variables**: camelCase for local variables, CamelCase for exported
- **Constants**: CamelCase or ALL_CAPS for package-level constants
- **Interfaces**: End with `-er` when possible (e.g., `DatabaseAnalyzer`)

### Code Organization

```go
package example

import (
    // Standard library imports first
    "context"
    "fmt"
    "time"
    
    // Third-party imports
    "github.com/some/package"
    
    // Local imports
    "github.com/cherry-pick/pkg/types"
)

// Constants
const DefaultTimeout = 30 * time.Second

// Types
type Service struct {
    connector interfaces.DatabaseConnector
    config    *types.Config
}

// Constructor
func NewService(connector interfaces.DatabaseConnector) *Service {
    return &Service{
        connector: connector,
    }
}

// Methods
func (s *Service) AnalyzeDatabase() (*types.DatabaseReport, error) {
    // Implementation
}
```

### Interface Design

Follow dependency injection patterns:

```go
// Define interfaces in the interfaces/ package
type DatabaseAnalyzer interface {
    AnalyzeTables() ([]types.TableInfo, error)
    AnalyzeSchema() (*types.SchemaInfo, error)
}

// Implement in specific packages
type MySQLAnalyzer struct {
    connector interfaces.DatabaseConnector
}

func (a *MySQLAnalyzer) AnalyzeTables() ([]types.TableInfo, error) {
    // Implementation
}
```

### Error Handling

- Use explicit error handling
- Wrap errors with context using `fmt.Errorf`
- Don't panic in library code
- Log errors appropriately

```go
func (s *Service) AnalyzeDatabase() (*types.DatabaseReport, error) {
    tables, err := s.analyzer.AnalyzeTables()
    if err != nil {
        return nil, fmt.Errorf("failed to analyze tables: %w", err)
    }
    
    // Continue processing
    return report, nil
}
```

### Documentation

- Document all exported functions, types, and constants
- Use complete sentences in comments
- Include examples for complex functions

```go
// AnalyzeDatabase performs comprehensive database analysis including
// table structure, data quality, and performance metrics.
//
// It returns a DatabaseReport containing all analysis results or an error
// if the analysis fails.
//
// Example:
//   report, err := service.AnalyzeDatabase()
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Health Score: %.2f", report.Summary.HealthScore)
func (s *Service) AnalyzeDatabase() (*types.DatabaseReport, error) {
    // Implementation
}
```

## Testing

### Test Structure

- Unit tests: `*_test.go` files alongside source code
- Mock interfaces for testing

### Writing Tests

```go
package analyzer

import (
    "testing"
    "github.com/cherry-pick/pkg/types"
)

func TestAnalyzeTables(t *testing.T) {
    // Setup
    mockConnector := &MockDatabaseConnector{
        tables: []string{"users", "orders"},
    }
    analyzer := NewDatabaseAnalyzer(mockConnector)
    
    // Execute
    tables, err := analyzer.AnalyzeTables()
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if len(tables) != 2 {
        t.Errorf("Expected 2 tables, got %d", len(tables))
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/analyzer/

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

### Test Coverage

Maintain minimum 80% test coverage for new code:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Mock Generation

Use interfaces for dependency injection and create mocks:

```go
type MockDatabaseConnector struct {
    tables []string
    err    error
}

func (m *MockDatabaseConnector) GetTables() ([]string, error) {
    return m.tables, m.err
}
```

## Pull Request Process

### Before Submitting

1. **Update from upstream:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run tests:**
   ```bash
   go test ./...
   go vet ./...
   ```

3. **Format code:**
   ```bash
   go fmt ./...
   goimports -w .
   ```

4. **Lint code:**
   ```bash
   golangci-lint run
   ```

### Commit Guidelines

Use conventional commit messages:

```
feat: add MongoDB support for data lineage tracking

- Implement MongoDB-specific lineage analysis
- Add aggregation pipeline for dependency detection
- Include tests for new functionality

Closes #123
```

**Commit types:**
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `style:` Code style changes
- `refactor:` Code refactoring
- `test:` Test additions or changes
- `chore:` Build system or auxiliary tool changes

### Pull Request Template

When creating a PR, include:

1. **Description**: What does this PR do?
2. **Motivation**: Why is this change needed?
3. **Testing**: How was this tested?
4. **Breaking Changes**: Any breaking changes?
5. **Checklist**: 
   - [ ] Tests pass
   - [ ] Code is formatted
   - [ ] Documentation updated
   - [ ] Changelog updated (if needed)

### Review Process

1. At least one maintainer review required
2. All tests must pass
3. Code coverage must not decrease significantly
4. Documentation must be updated for new features

## Issue Reporting

### Bug Reports

Include:
- Go version (`go version`)
- Operating system
- Database type and version
- Steps to reproduce
- Expected vs actual behavior
- Error messages or logs

### Feature Requests

Include:
- Use case description
- Proposed solution
- Alternative solutions considered
- Additional context

### Issue Labels

- `bug`: Something isn't working
- `enhancement`: New feature or improvement
- `documentation`: Documentation improvements
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention needed

## Documentation

### Code Documentation

- All exported functions must have doc comments
- Include examples for complex APIs
- Document error conditions
- Explain non-obvious behavior

### README Updates

- Update feature lists for new capabilities
- Add new configuration options
- Include new examples
- Update supported database versions

### Wiki and Guides

- Create guides for complex features
- Add troubleshooting information
- Include architecture decisions
- Document best practices

## Development Workflow

### Feature Development

1. **Create feature branch:**
   ```bash
   git checkout -b feature/mongodb-optimization
   ```

2. **Implement feature:**
   - Write tests first (TDD approach recommended)
   - Implement functionality
   - Update documentation

3. **Test thoroughly:**
   ```bash
   go test ./...
   go run examples/env-database-example.go
   ```

4. **Submit PR:**
   - Push to your fork
   - Create pull request
   - Address review feedback

### Hotfix Process

1. **Create hotfix branch from main:**
   ```bash
   git checkout -b hotfix/security-fix
   ```

2. **Make minimal changes**
3. **Test thoroughly**
4. **Submit PR with `hotfix` label**

## Getting Help

- **GitHub Discussions**: For questions and general discussion
- **GitHub Issues**: For bug reports and feature requests
- **Code Review**: For feedback on implementation approaches

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes for significant contributions
- GitHub contributor graphs

Thank you for contributing to the Enterprise Database Intelligence System!
