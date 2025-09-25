package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cherry-pick/pkg/analyzer/core"
)

type DatabaseAnalyzerService struct {
	connector   core.DatabaseConnector
	calculator  core.AnalysisCalculator
	aggregator  core.AnalysisAggregator
	validator   core.AnalysisValidator
}

func NewDatabaseAnalyzerService(
	connector core.DatabaseConnector,
	calculator core.AnalysisCalculator,
	aggregator core.AnalysisAggregator,
	validator core.AnalysisValidator,
) *DatabaseAnalyzerService {
	return &DatabaseAnalyzerService{
		connector:  connector,
		calculator: calculator,
		aggregator: aggregator,
		validator:  validator,
	}
}

func (das *DatabaseAnalyzerService) AnalyzeDatabase(ctx context.Context, request core.AnalysisRequest) (*core.AnalysisResult, error) {
	if err := das.validator.ValidateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid analysis request: %w", err)
	}

	if !das.connector.IsConnected() {
		return nil, fmt.Errorf("database not connected")
	}

	startTime := time.Now()
	log.Printf("Starting database analysis for %s", request.DatabaseType)

	tables, err := das.AnalyzeTables(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze tables: %w", err)
	}

	summary := das.aggregator.AggregateTableStats(tables)
	insights := das.generateInsights(tables)
	recommendations := das.generateRecommendations(insights)

	var performance *core.PerformanceMetrics
	if request.Options.IncludePerformance {
		performance, err = das.GetPerformanceMetrics(ctx, request)
		if err != nil {
			log.Printf("Warning: Could not get performance metrics: %v", err)
		}
	}

	result := &core.AnalysisResult{
		ID:             generateAnalysisID(),
		DatabaseName:   das.connector.GetDatabaseName(),
		DatabaseType:   request.DatabaseType,
		AnalysisTime:   time.Now(),
		Summary:        summary,
		Tables:         tables,
		Insights:       insights,
		Recommendations: recommendations,
		Performance:    performance,
	}

	log.Printf("Database analysis completed in %v", time.Since(startTime))
	return result, nil
}

func (das *DatabaseAnalyzerService) AnalyzeTables(ctx context.Context, request core.AnalysisRequest) ([]core.TableInfo, error) {
	tableNames, err := das.GetTableNames(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get table names: %w", err)
	}

	var tables []core.TableInfo
	for _, tableName := range tableNames {
		log.Printf("Analyzing table: %s", tableName)

		table, err := das.AnalyzeTable(ctx, tableName, request)
		if err != nil {
			log.Printf("Warning: Failed to analyze table %s: %v", tableName, err)
			continue
		}
		tables = append(tables, *table)
	}

	return tables, nil
}

func (das *DatabaseAnalyzerService) AnalyzeTable(ctx context.Context, tableName string, request core.AnalysisRequest) (*core.TableInfo, error) {
	table := &core.TableInfo{
		Name:         tableName,
		LastModified: time.Now(),
	}

	if request.Options.IncludeData {
		rowCount, err := das.getRowCount(ctx, tableName)
		if err != nil {
			log.Printf("Warning: Could not get row count for %s: %v", tableName, err)
		}
		table.RowCount = rowCount

		size, err := das.getTableSize(ctx, tableName)
		if err != nil {
			log.Printf("Warning: Could not get table size for %s: %v", tableName, err)
		}
		table.Size = size
	}

	if request.Options.IncludeSchema {
		columns, err := das.analyzeColumns(ctx, tableName, request)
		if err != nil {
			return table, fmt.Errorf("failed to analyze columns: %w", err)
		}
		table.Columns = columns
	}

	if request.Options.IncludeIndexes {
		indexes, err := das.getIndexes(ctx, tableName)
		if err != nil {
			log.Printf("Warning: Could not get indexes for %s: %v", tableName, err)
		}
		table.Indexes = indexes
	}

	if request.Options.IncludeRelations {
		constraints, err := das.getConstraints(ctx, tableName)
		if err != nil {
			log.Printf("Warning: Could not get constraints for %s: %v", tableName, err)
		}
		table.Constraints = constraints

		relationships, err := das.getRelationships(ctx, tableName)
		if err != nil {
			log.Printf("Warning: Could not get relationships for %s: %v", tableName, err)
		}
		table.Relationships = relationships
	}

	return table, nil
}

func (das *DatabaseAnalyzerService) GetTableNames(ctx context.Context, request core.AnalysisRequest) ([]string, error) {
	db := das.connector.GetDatabase().(*sql.DB)
	dbType := string(request.DatabaseType)

	var query string
	switch dbType {
	case "mysql":
		query = "SHOW TABLES"
	case "postgres":
		query = "SELECT tablename FROM pg_tables WHERE schemaname = 'public'"
	case "sqlite3":
		query = "SELECT name FROM sqlite_master WHERE type='table'"
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table names: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	return tables, rows.Err()
}

func (das *DatabaseAnalyzerService) GetTableStats(ctx context.Context, tableName string, request core.AnalysisRequest) (*core.TableInfo, error) {
	return das.AnalyzeTable(ctx, tableName, request)
}

func (das *DatabaseAnalyzerService) GetPerformanceMetrics(ctx context.Context, request core.AnalysisRequest) (*core.PerformanceMetrics, error) {
	db := das.connector.GetDatabase().(*sql.DB)
	dbType := string(request.DatabaseType)

	metrics := &core.PerformanceMetrics{}

	switch dbType {
	case "mysql":
		return das.getMySQLPerformanceMetrics(ctx, db)
	case "postgres":
		return das.getPostgresPerformanceMetrics(ctx, db)
	case "sqlite3":
		return das.getSQLitePerformanceMetrics(ctx, db)
	default:
		return metrics, fmt.Errorf("performance metrics not supported for database type: %s", dbType)
	}
}

func (das *DatabaseAnalyzerService) getRowCount(ctx context.Context, tableName string) (int64, error) {
	db := das.connector.GetDatabase().(*sql.DB)
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	
	var count int64
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get row count: %w", err)
	}
	return count, nil
}

func (das *DatabaseAnalyzerService) getTableSize(ctx context.Context, tableName string) (string, error) {
	db := das.connector.GetDatabase().(*sql.DB)
	dbType := das.connector.GetDatabaseType()

	var query string
	var args []interface{}

	switch dbType {
	case core.DatabaseTypeMySQL:
		query = `
			SELECT 
				ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
			FROM information_schema.tables 
			WHERE table_schema = DATABASE() AND table_name = ?`
		args = append(args, tableName)
	case core.DatabaseTypePostgres:
		query = `SELECT pg_size_pretty(pg_total_relation_size($1)) AS size`
		args = append(args, tableName)
	case core.DatabaseTypeSQLite:
		query = `SELECT COUNT(*) * 1024 as approx_bytes FROM pragma_table_info(?) LIMIT 1`
		args = append(args, tableName)
	default:
		return "Unknown", fmt.Errorf("unsupported database type: %s", dbType)
	}

	var sizeResult interface{}
	err := db.QueryRowContext(ctx, query, args...).Scan(&sizeResult)
	if err != nil {
		return "Unknown", fmt.Errorf("failed to get table size: %w", err)
	}

	switch dbType {
	case core.DatabaseTypeMySQL:
		if size, ok := sizeResult.(float64); ok {
			return fmt.Sprintf("%.2f MB", size), nil
		}
	case core.DatabaseTypePostgres:
		if size, ok := sizeResult.(string); ok {
			return size, nil
		}
	case core.DatabaseTypeSQLite:
		if size, ok := sizeResult.(int64); ok {
			if size < 1024 {
				return fmt.Sprintf("%d bytes", size), nil
			} else if size < 1024*1024 {
				return fmt.Sprintf("%.2f KB", float64(size)/1024), nil
			} else {
				return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024)), nil
			}
		}
	}

	return "Unknown", nil
}

func (das *DatabaseAnalyzerService) analyzeColumns(ctx context.Context, tableName string, request core.AnalysisRequest) ([]core.ColumnInfo, error) {
	db := das.connector.GetDatabase().(*sql.DB)
	dbType := das.connector.GetDatabaseType()

	var query string
	var args []interface{}

	switch dbType {
	case core.DatabaseTypeMySQL:
		query = `
			SELECT 
				COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT,
				CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE,
				COLUMN_KEY
			FROM INFORMATION_SCHEMA.COLUMNS 
			WHERE TABLE_NAME = ? AND TABLE_SCHEMA = DATABASE()
			ORDER BY ORDINAL_POSITION`
		args = append(args, tableName)
	case core.DatabaseTypePostgres:
		query = `
			SELECT 
				column_name, data_type, is_nullable, column_default,
				character_maximum_length, numeric_precision, numeric_scale,
				CASE WHEN column_name IN (
					SELECT column_name FROM information_schema.table_constraints tc
					JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
					WHERE tc.table_name = $1 AND tc.constraint_type = 'PRIMARY KEY'
				) THEN 'PRI' ELSE '' END as column_key
			FROM information_schema.columns 
			WHERE table_name = $1
			ORDER BY ordinal_position`
		args = append(args, tableName)
	case core.DatabaseTypeSQLite:
		query = fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []core.ColumnInfo
	for rows.Next() {
		var col core.ColumnInfo
		var maxLength, precision, scale sql.NullInt64
		var defaultVal sql.NullString
		var columnKey string

		if dbType == core.DatabaseTypeSQLite {
			var cid int
			var notNull int
			var pk int
			err = rows.Scan(&cid, &col.Name, &col.DataType, &notNull, &defaultVal, &pk)
			col.IsNullable = notNull == 0
			col.IsPrimaryKey = pk == 1
		} else {
			var nullable string
			err = rows.Scan(&col.Name, &col.DataType, &nullable, &defaultVal,
				&maxLength, &precision, &scale, &columnKey)
			col.IsNullable = nullable == "YES"
			col.IsPrimaryKey = columnKey == "PRI"
		}

		if err != nil {
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}

		if defaultVal.Valid {
			col.DefaultValue = defaultVal.String
		}
		if maxLength.Valid {
			col.MaxLength = int(maxLength.Int64)
		}
		if precision.Valid {
			col.Precision = int(precision.Int64)
		}
		if scale.Valid {
			col.Scale = int(scale.Int64)
		}

		if request.Options.IncludeData {
			col.DataProfile = das.analyzeColumnData(ctx, tableName, col.Name, col.DataType)
			col.UniqueValues = das.getUniqueValueCount(ctx, tableName, col.Name)
			col.NullCount = das.getNullCount(ctx, tableName, col.Name)
		}

		columns = append(columns, col)
	}

	return columns, rows.Err()
}

func (das *DatabaseAnalyzerService) analyzeColumnData(ctx context.Context, tableName, columnName, dataType string) core.DataProfile {
	profile := core.DataProfile{}

	db := das.connector.GetDatabase().(*sql.DB)
	sampleQuery := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s IS NOT NULL LIMIT 10",
		columnName, tableName, columnName)

	rows, err := db.QueryContext(ctx, sampleQuery)
	if err != nil {
		return profile
	}
	defer rows.Close()

	var samples []string
	for rows.Next() {
		var value sql.NullString
		if err := rows.Scan(&value); err != nil {
			continue
		}
		if value.Valid {
			samples = append(samples, value.String)
		}
	}
	profile.SampleData = samples

	if das.isNumericType(dataType) {
		minQuery := fmt.Sprintf("SELECT MIN(%s), MAX(%s), AVG(%s) FROM %s WHERE %s IS NOT NULL",
			columnName, columnName, columnName, tableName, columnName)

		var min, max, avg sql.NullFloat64
		err := db.QueryRowContext(ctx, minQuery).Scan(&min, &max, &avg)
		if err == nil {
			if min.Valid {
				profile.Min = min.Float64
			}
			if max.Valid {
				profile.Max = max.Float64
			}
			if avg.Valid {
				profile.Avg = avg.Float64
			}
		}
	}

	profile.Quality = das.calculator.CalculateDataQuality(core.ColumnInfo{
		Name: columnName,
		DataProfile: profile,
	})

	return profile
}

func (das *DatabaseAnalyzerService) getUniqueValueCount(ctx context.Context, tableName, columnName string) int64 {
	db := das.connector.GetDatabase().(*sql.DB)
	query := fmt.Sprintf("SELECT COUNT(DISTINCT %s) FROM %s", columnName, tableName)
	
	var count int64
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (das *DatabaseAnalyzerService) getNullCount(ctx context.Context, tableName, columnName string) int64 {
	db := das.connector.GetDatabase().(*sql.DB)
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s IS NULL", tableName, columnName)
	
	var count int64
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (das *DatabaseAnalyzerService) getIndexes(ctx context.Context, tableName string) ([]core.IndexInfo, error) {
	db := das.connector.GetDatabase().(*sql.DB)
	dbType := das.connector.GetDatabaseType()

	var query string
	var args []interface{}

	switch dbType {
	case core.DatabaseTypeMySQL:
		query = "SHOW INDEX FROM " + tableName
	case core.DatabaseTypePostgres:
		query = `
			SELECT 
				indexname, 
				indexdef,
				CASE WHEN indisunique THEN true ELSE false END as is_unique
			FROM pg_indexes 
			JOIN pg_class ON pg_class.relname = indexname
			JOIN pg_index ON pg_index.indexrelid = pg_class.oid
			WHERE tablename = $1`
		args = append(args, tableName)
	case core.DatabaseTypeSQLite:
		query = "PRAGMA index_list(" + tableName + ")"
	default:
		return []core.IndexInfo{}, fmt.Errorf("unsupported database type: %s", dbType)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer rows.Close()

	var indexes []core.IndexInfo
	for rows.Next() {
		var index core.IndexInfo

		switch dbType {
		case core.DatabaseTypeMySQL:
			var seq, cardinality sql.NullInt64
			var collation, subPart, packed, null, comment sql.NullString
			var nonUnique int
			var columnName string
			err = rows.Scan(&tableName, &nonUnique, &index.Name, &seq,
				&columnName, &collation, &cardinality, &subPart, &packed, &null, &index.Type, &comment)
			if err != nil {
				continue
			}
			index.IsUnique = nonUnique == 0
			index.Columns = []string{columnName}
		case core.DatabaseTypePostgres:
			var indexDef string
			err = rows.Scan(&index.Name, &indexDef, &index.IsUnique)
			if err != nil {
				continue
			}
			index.Type = "btree"
			index.Columns = []string{"parsed_from_def"}
		case core.DatabaseTypeSQLite:
			var seq int
			var unique int
			var origin string
			var partial int
			err = rows.Scan(&seq, &index.Name, &unique, &origin, &partial)
			if err != nil {
				continue
			}
			index.IsUnique = unique == 1
			index.Type = "btree"

			colQuery := "PRAGMA index_info(" + index.Name + ")"
			colRows, colErr := db.QueryContext(ctx, colQuery)
			if colErr == nil {
				var columns []string
				for colRows.Next() {
					var seqno, cid int
					var name string
					if colRows.Scan(&seqno, &cid, &name) == nil {
						columns = append(columns, name)
					}
				}
				colRows.Close()
				index.Columns = columns
			}
		}

		if index.Name != "" {
			indexes = append(indexes, index)
		}
	}

	return indexes, rows.Err()
}

func (das *DatabaseAnalyzerService) getConstraints(ctx context.Context, tableName string) ([]core.Constraint, error) {
	db := das.connector.GetDatabase().(*sql.DB)
	dbType := das.connector.GetDatabaseType()

	var query string
	var args []interface{}

	switch dbType {
	case core.DatabaseTypeMySQL:
		query = `
			SELECT 
				CONSTRAINT_NAME, CONSTRAINT_TYPE, COLUMN_NAME
			FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
			LEFT JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu 
				ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME 
				AND tc.TABLE_NAME = kcu.TABLE_NAME
			WHERE tc.TABLE_NAME = ? AND tc.TABLE_SCHEMA = DATABASE()`
		args = append(args, tableName)
	case core.DatabaseTypePostgres:
		query = `
			SELECT 
				tc.constraint_name, tc.constraint_type, kcu.column_name
			FROM information_schema.table_constraints tc
			LEFT JOIN information_schema.key_column_usage kcu 
				ON tc.constraint_name = kcu.constraint_name 
				AND tc.table_name = kcu.table_name
			WHERE tc.table_name = $1`
		args = append(args, tableName)
	case core.DatabaseTypeSQLite:
		query = "PRAGMA foreign_key_list(" + tableName + ")"
	default:
		return []core.Constraint{}, fmt.Errorf("unsupported database type: %s", dbType)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query constraints: %w", err)
	}
	defer rows.Close()

	var constraints []core.Constraint
	for rows.Next() {
		var constraint core.Constraint

		switch dbType {
		case core.DatabaseTypeMySQL, core.DatabaseTypePostgres:
			var columnName sql.NullString
			err = rows.Scan(&constraint.Name, &constraint.Type, &columnName)
			if err != nil {
				continue
			}
			if columnName.Valid {
				constraint.Columns = []string{columnName.String}
			}
		case core.DatabaseTypeSQLite:
			var id int
			var seq int
			var table, from, to, onUpdate, onDelete, match string
			err = rows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match)
			if err != nil {
				continue
			}
			constraint.Name = fmt.Sprintf("fk_%d", id)
			constraint.Type = "FOREIGN KEY"
			constraint.Columns = []string{from}
			constraint.RefTable = table
			constraint.RefColumns = []string{to}
		}

		if constraint.Name != "" {
			constraints = append(constraints, constraint)
		}
	}

	return constraints, rows.Err()
}

func (das *DatabaseAnalyzerService) getRelationships(ctx context.Context, tableName string) ([]core.Relationship, error) {
	db := das.connector.GetDatabase().(*sql.DB)
	dbType := das.connector.GetDatabaseType()

	var query string
	var args []interface{}

	switch dbType {
	case core.DatabaseTypeMySQL:
		query = `
			SELECT 
				kcu.COLUMN_NAME,
				kcu.REFERENCED_TABLE_NAME,
				kcu.REFERENCED_COLUMN_NAME,
				rc.UPDATE_RULE,
				rc.DELETE_RULE
			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc 
				ON kcu.CONSTRAINT_NAME = rc.CONSTRAINT_NAME
			WHERE kcu.TABLE_NAME = ? 
				AND kcu.REFERENCED_TABLE_NAME IS NOT NULL
				AND kcu.TABLE_SCHEMA = DATABASE()`
		args = append(args, tableName)
	case core.DatabaseTypePostgres:
		query = `
			SELECT 
				kcu.column_name,
				ccu.table_name AS referenced_table,
				ccu.column_name AS referenced_column,
				rc.update_rule,
				rc.delete_rule
			FROM information_schema.key_column_usage kcu
			JOIN information_schema.referential_constraints rc 
				ON kcu.constraint_name = rc.constraint_name
			JOIN information_schema.constraint_column_usage ccu 
				ON rc.unique_constraint_name = ccu.constraint_name
			WHERE kcu.table_name = $1`
		args = append(args, tableName)
	case core.DatabaseTypeSQLite:
		query = "PRAGMA foreign_key_list(" + tableName + ")"
	default:
		return []core.Relationship{}, fmt.Errorf("unsupported database type: %s", dbType)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships: %w", err)
	}
	defer rows.Close()

	var relationships []core.Relationship
	for rows.Next() {
		var rel core.Relationship

		switch dbType {
		case core.DatabaseTypeMySQL, core.DatabaseTypePostgres:
			var updateRule, deleteRule sql.NullString
			err = rows.Scan(&rel.SourceColumn, &rel.TargetTable,
				&rel.TargetColumn, &updateRule, &deleteRule)
			if err != nil {
				continue
			}
			rel.Type = "FOREIGN KEY"
		case core.DatabaseTypeSQLite:
			var id, seq int
			var onUpdate, onDelete, match string
			err = rows.Scan(&id, &seq, &rel.TargetTable, &rel.SourceColumn,
				&rel.TargetColumn, &onUpdate, &onDelete, &match)
			if err != nil {
				continue
			}
			rel.Type = "FOREIGN KEY"
		}

		if rel.TargetTable != "" {
			relationships = append(relationships, rel)
		}
	}

	return relationships, rows.Err()
}

func (das *DatabaseAnalyzerService) getMySQLPerformanceMetrics(ctx context.Context, db *sql.DB) (*core.PerformanceMetrics, error) {
	metrics := &core.PerformanceMetrics{}

	query := "SHOW STATUS LIKE 'Connections'"
	var name string
	var value int
	err := db.QueryRowContext(ctx, query).Scan(&name, &value)
	if err == nil {
		metrics.Connections.TotalCreated = value
	}

	return metrics, nil
}

func (das *DatabaseAnalyzerService) getPostgresPerformanceMetrics(ctx context.Context, db *sql.DB) (*core.PerformanceMetrics, error) {
	metrics := &core.PerformanceMetrics{}

	query := "SELECT count(*) FROM pg_stat_activity"
	var connections int
	err := db.QueryRowContext(ctx, query).Scan(&connections)
	if err == nil {
		metrics.Connections.Current = connections
	}

	return metrics, nil
}

func (das *DatabaseAnalyzerService) getSQLitePerformanceMetrics(ctx context.Context, db *sql.DB) (*core.PerformanceMetrics, error) {
	return &core.PerformanceMetrics{}, nil
}

func (das *DatabaseAnalyzerService) isNumericType(dataType string) bool {
	numericTypes := []string{"int", "integer", "bigint", "smallint", "tinyint", "decimal", "numeric", "float", "double", "real"}
	for _, t := range numericTypes {
		if dataType == t {
			return true
		}
	}
	return false
}

func (das *DatabaseAnalyzerService) generateInsights(tables []core.TableInfo) []core.DatabaseInsight {
	var insights []core.DatabaseInsight

	for _, table := range tables {
		if table.RowCount > 1000000 {
			insight := core.DatabaseInsight{
				Type:           "performance",
				Severity:       "medium",
				Title:          "Large Table Detected",
				Description:    fmt.Sprintf("Table '%s' has %d rows, which may impact performance", table.Name, table.RowCount),
				Suggestion:     "Consider partitioning, archiving old data, or optimizing queries",
				AffectedTables: []string{table.Name},
				MetricValue:    table.RowCount,
			}
			insights = append(insights, insight)
		}

		if len(table.Indexes) <= 1 && table.RowCount > 10000 {
			insight := core.DatabaseInsight{
				Type:           "performance",
				Severity:       "high",
				Title:          "Missing Indexes on Large Table",
				Description:    fmt.Sprintf("Table '%s' has %d rows but only %d indexes", table.Name, table.RowCount, len(table.Indexes)),
				Suggestion:     "Consider adding indexes on frequently queried columns",
				AffectedTables: []string{table.Name},
				MetricValue:    len(table.Indexes),
			}
			insights = append(insights, insight)
		}
	}

	return insights
}

func (das *DatabaseAnalyzerService) generateRecommendations(insights []core.DatabaseInsight) []string {
	var recommendations []string

	highPriorityCount := 0
	for _, insight := range insights {
		if insight.Severity == "high" {
			highPriorityCount++
			recommendations = append(recommendations, fmt.Sprintf("Priority: %s", insight.Suggestion))
		}
	}

	if highPriorityCount == 0 {
		recommendations = append(recommendations, "Database appears to be in good condition with no critical issues")
	}

	return recommendations
}

func generateAnalysisID() string {
	return fmt.Sprintf("analysis_%d", time.Now().Unix())
}
