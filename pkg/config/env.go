package config

import (
	"fmt"
	"os"
	"strings"
)

type DatabaseConfig struct {
	URL      string
	Type     string
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	config := &DatabaseConfig{}

	if url := os.Getenv("DATABASE_URL"); url != "" {
		config.URL = url
		config.Type = detectDatabaseType(url)
		return config, nil
	}

	if mongoURL := os.Getenv("MONGODB_URL"); mongoURL != "" {
		config.URL = mongoURL
		config.Type = "mongodb"
		return config, nil
	}

	config.Type = strings.ToLower(os.Getenv("DB_TYPE"))
	config.Host = getEnvOrDefault("DB_HOST", "localhost")
	config.Port = os.Getenv("DB_PORT")
	config.Database = os.Getenv("DB_NAME")
	config.Username = os.Getenv("DB_USER")
	config.Password = os.Getenv("DB_PASSWORD")
	config.SSLMode = getEnvOrDefault("DB_SSLMODE", "disable")

	if config.Type != "" && config.Database != "" {
		var err error
		config.URL, err = buildConnectionString(config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	return nil, fmt.Errorf("no database configuration found in environment variables")
}

func GetAllSupportedConfigs() map[string][]string {
	return map[string][]string{
		"Complete URLs": {
			`DATABASE_URL="mysql://user:password@localhost:3306/dbname"`,
			`DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"`,
			`DATABASE_URL="sqlite3://./database.db"`,
			`MONGODB_URL="mongodb://localhost:27017/dbname"`,
			`MONGODB_URL="mongodb+srv://user:password@cluster.mongodb.net/dbname"`,
		},
		"MySQL Components": {
			`DB_TYPE="mysql"`,
			`DB_HOST="localhost"`,
			`DB_PORT="3306"`,
			`DB_NAME="mydb"`,
			`DB_USER="root"`,
			`DB_PASSWORD="password"`,
		},
		"PostgreSQL Components": {
			`DB_TYPE="postgres"`,
			`DB_HOST="localhost"`,
			`DB_PORT="5432"`,
			`DB_NAME="mydb"`,
			`DB_USER="postgres"`,
			`DB_PASSWORD="password"`,
			`DB_SSLMODE="disable"`,
		},
		"SQLite Components": {
			`DB_TYPE="sqlite3"`,
			`DB_NAME="./database.db"`,
		},
	}
}

func PrintConfigHelp() {
	fmt.Println("=== Database Configuration Help ===")
	fmt.Println("You can configure database connections using environment variables in several ways:")
	fmt.Println()

	configs := GetAllSupportedConfigs()
	for category, examples := range configs {
		fmt.Printf("%s:\n", category)
		for _, example := range examples {
			fmt.Printf("  %s\n", example)
		}
		fmt.Println()
	}

	fmt.Println("Examples:")
	fmt.Println("  # Using complete URL")
	fmt.Println("  $env:DATABASE_URL = \"mysql://root:password@localhost:3306/testdb\"")
	fmt.Println("  go run examples/real-database-example.go")
	fmt.Println()
	fmt.Println("  # Using individual components")
	fmt.Println("  $env:DB_TYPE = \"postgres\"")
	fmt.Println("  $env:DB_HOST = \"localhost\"")
	fmt.Println("  $env:DB_NAME = \"mydb\"")
	fmt.Println("  $env:DB_USER = \"postgres\"")
	fmt.Println("  $env:DB_PASSWORD = \"password\"")
	fmt.Println("  go run examples/real-database-example.go")
}

func detectDatabaseType(url string) string {
	url = strings.ToLower(url)
	switch {
	case strings.HasPrefix(url, "mysql://") || strings.Contains(url, "@tcp"):
		return "mysql"
	case strings.HasPrefix(url, "postgres://") || strings.HasPrefix(url, "postgresql://"):
		return "postgres"
	case strings.HasPrefix(url, "sqlite3://") || strings.HasSuffix(url, ".db"):
		return "sqlite3"
	case strings.HasPrefix(url, "mongodb://") || strings.HasPrefix(url, "mongodb+srv://"):
		return "mongodb"
	default:
		return "unknown"
	}
}

func buildConnectionString(config *DatabaseConfig) (string, error) {
	switch config.Type {
	case "mysql":
		port := getEnvOrDefault("DB_PORT", "3306")
		if config.Username != "" && config.Password != "" {
			return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
				config.Username, config.Password, config.Host, port, config.Database), nil
		}
		return fmt.Sprintf("tcp(%s:%s)/%s", config.Host, port, config.Database), nil

	case "postgres", "postgresql":
		port := getEnvOrDefault("DB_PORT", "5432")
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.Username, config.Password, config.Host, port, config.Database, config.SSLMode), nil

	case "sqlite3":
		return config.Database, nil

	case "mongodb":
		port := getEnvOrDefault("DB_PORT", "27017")
		if config.Username != "" && config.Password != "" {
			return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
				config.Username, config.Password, config.Host, port, config.Database), nil
		}
		return fmt.Sprintf("mongodb://%s:%s/%s", config.Host, port, config.Database), nil

	default:
		return "", fmt.Errorf("unsupported database type: %s", config.Type)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
