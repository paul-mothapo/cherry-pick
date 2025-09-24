package connector

import (
	"database/sql"
	"fmt"

	"github.com/cherry-pick/pkg/interfaces"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseConnectorImpl struct {
	db             *sql.DB
	driverName     string
	dataSourceName string
}

func NewDatabaseConnector(driverName, dataSourceName string) interfaces.DatabaseConnector {
	return &DatabaseConnectorImpl{
		driverName:     driverName,
		dataSourceName: dataSourceName,
	}
}

func (dc *DatabaseConnectorImpl) Connect() error {
	db, err := sql.Open(dc.driverName, dc.dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	dc.db = db
	return nil
}

func (dc *DatabaseConnectorImpl) Close() error {
	if dc.db == nil {
		return nil
	}

	err := dc.db.Close()
	dc.db = nil
	return err
}

func (dc *DatabaseConnectorImpl) Ping() error {
	if dc.db == nil {
		return fmt.Errorf("database connection is not established")
	}
	return dc.db.Ping()
}

func (dc *DatabaseConnectorImpl) GetDB() *sql.DB {
	return dc.db
}

func (dc *DatabaseConnectorImpl) GetDatabaseName() (string, error) {
	if dc.db == nil {
		return "", fmt.Errorf("database connection is not established")
	}

	var query string
	switch dc.driverName {
	case "mysql":
		query = "SELECT DATABASE()"
	case "postgres":
		query = "SELECT current_database()"
	case "sqlite3":
		return "SQLite Database", nil
	default:
		return "Unknown", nil
	}

	var name sql.NullString
	err := dc.db.QueryRow(query).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("failed to get database name: %w", err)
	}

	if name.Valid {
		return name.String, nil
	}
	return "Unknown", nil
}

func (dc *DatabaseConnectorImpl) GetDatabaseType() string {
	return dc.driverName
}
