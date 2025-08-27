package analyzer

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/intelligent-algorithm/pkg/interfaces"
	"github.com/intelligent-algorithm/pkg/types"
	"github.com/intelligent-algorithm/pkg/utils"
)

type DatabaseAnalyzerImpl struct {
	db     *sql.DB
	dbType string
}

func NewDatabaseAnalyzer(db *sql.DB, dbType string) interfaces.DatabaseAnalyzer {
	return &DatabaseAnalyzerImpl{
		db:     db,
		dbType: dbType,
	}
}

func (da *DatabaseAnalyzerImpl) AnalyzeTables() ([]types.TableInfo, error) {
	tableNames, err := da.GetTableNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get table names: %w", err)
	}

	var tables []types.TableInfo
	for _, tableName := range tableNames {
		log.Printf("Analyzing table: %s", tableName)

		tableInfo, err := da.AnalyzeTable(tableName)
		if err != nil {
			log.Printf("Warning: Failed to analyze table %s: %v", tableName, err)
			continue
		}
		tables = append(tables, tableInfo)
	}

	return tables, nil
}

func (da *DatabaseAnalyzerImpl) AnalyzeTable(tableName string) (types.TableInfo, error) {
	tableInfo := types.TableInfo{Name: tableName}

	rowCount, err := da.GetRowCount(tableName)
	if err != nil {
		log.Printf("Warning: Could not get row count for %s: %v", tableName, err)
		rowCount = 0
	}
	tableInfo.RowCount = rowCount

	columns, err := da.AnalyzeColumns(tableName)
	if err != nil {
		return tableInfo, fmt.Errorf("failed to analyze columns: %w", err)
	}
	tableInfo.Columns = columns

	indexes, err := da.GetIndexes(tableName)
	if err != nil {
		log.Printf("Warning: Could not get indexes for %s: %v", tableName, err)
	}
	tableInfo.Indexes = indexes

	constraints, err := da.GetConstraints(tableName)
	if err != nil {
		log.Printf("Warning: Could not get constraints for %s: %v", tableName, err)
	}
	tableInfo.Constraints = constraints

	relationships, err := da.GetRelationships(tableName)
	if err != nil {
		log.Printf("Warning: Could not get relationships for %s: %v", tableName, err)
	}
	tableInfo.Relationships = relationships

	size, err := da.GetTableSize(tableName)
	if err != nil {
		log.Printf("Warning: Could not get table size for %s: %v", tableName, err)
		size = "Unknown"
	}
	tableInfo.Size = size

	return tableInfo, nil
}

func (da *DatabaseAnalyzerImpl) AnalyzeColumns(tableName string) ([]types.ColumnInfo, error) {
	var query string
	switch da.dbType {
	case "mysql":
		query = `
			SELECT 
				COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT,
				CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE,
				COLUMN_KEY
			FROM INFORMATION_SCHEMA.COLUMNS 
			WHERE TABLE_NAME = ? AND TABLE_SCHEMA = DATABASE()
			ORDER BY ORDINAL_POSITION`
	case "postgres":
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
	case "sqlite3":
		query = fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", da.dbType)
	}

	rows, err := da.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []types.ColumnInfo
	for rows.Next() {
		var col types.ColumnInfo
		var maxLength, precision, scale sql.NullInt64
		var defaultVal sql.NullString
		var columnKey string

		if da.dbType == "sqlite3" {
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

		profile, err := da.AnalyzeColumnData(tableName, col.Name, col.DataType)
		if err != nil {
			log.Printf("Warning: Could not analyze column data for %s.%s: %v",
				tableName, col.Name, err)
		}
		col.DataProfile = profile

		uniqueCount, err := da.GetUniqueValueCount(tableName, col.Name)
		if err != nil {
			log.Printf("Warning: Could not get unique count for %s.%s: %v",
				tableName, col.Name, err)
		}
		col.UniqueValues = uniqueCount

		nullCount, err := da.GetNullCount(tableName, col.Name)
		if err != nil {
			log.Printf("Warning: Could not get null count for %s.%s: %v",
				tableName, col.Name, err)
		}
		col.NullCount = nullCount

		columns = append(columns, col)
	}

	return columns, rows.Err()
}

func (da *DatabaseAnalyzerImpl) AnalyzeColumnData(tableName, columnName, dataType string) (types.DataProfile, error) {
	profile := types.DataProfile{}

	sampleQuery := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s IS NOT NULL LIMIT 10",
		columnName, tableName, columnName)

	rows, err := da.db.Query(sampleQuery)
	if err != nil {
		return profile, fmt.Errorf("failed to get sample data: %w", err)
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

	if utils.IsNumericType(dataType) {
		minQuery := fmt.Sprintf("SELECT MIN(%s), MAX(%s), AVG(%s) FROM %s WHERE %s IS NOT NULL",
			columnName, columnName, columnName, tableName, columnName)

		var min, max, avg sql.NullFloat64
		err := da.db.QueryRow(minQuery).Scan(&min, &max, &avg)
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

	if utils.IsStringType(dataType) && len(samples) > 0 {
		profile.Pattern = utils.DetectPattern(samples)
	}

	profile.Quality = da.CalculateDataQuality(tableName, columnName)

	return profile, nil
}

func (da *DatabaseAnalyzerImpl) GetTableNames() ([]string, error) {
	var query string
	switch da.dbType {
	case "mysql":
		query = "SHOW TABLES"
	case "postgres":
		query = "SELECT tablename FROM pg_tables WHERE schemaname = 'public'"
	case "sqlite3":
		query = "SELECT name FROM sqlite_master WHERE type='table'"
	default:
		return nil, fmt.Errorf("unsupported database type: %s", da.dbType)
	}

	rows, err := da.db.Query(query)
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

func (da *DatabaseAnalyzerImpl) GetRowCount(tableName string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var count int64
	err := da.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get row count: %w", err)
	}
	return count, nil
}

func (da *DatabaseAnalyzerImpl) GetUniqueValueCount(tableName, columnName string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(DISTINCT %s) FROM %s", columnName, tableName)
	var count int64
	err := da.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unique value count: %w", err)
	}
	return count, nil
}

func (da *DatabaseAnalyzerImpl) GetNullCount(tableName, columnName string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s IS NULL", tableName, columnName)
	var count int64
	err := da.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get null count: %w", err)
	}
	return count, nil
}

func (da *DatabaseAnalyzerImpl) CalculateDataQuality(tableName, columnName string) float64 {
	totalRows, err := da.GetRowCount(tableName)
	if err != nil || totalRows == 0 {
		return 1.0
	}

	nullCount, err := da.GetNullCount(tableName, columnName)
	if err != nil {
		return 1.0
	}

	nullRatio := float64(nullCount) / float64(totalRows)
	return 1.0 - nullRatio
}

func (da *DatabaseAnalyzerImpl) GetIndexes(tableName string) ([]types.IndexInfo, error) {
	var query string
	var args []interface{}

	switch da.dbType {
	case "mysql":
		query = "SHOW INDEX FROM " + tableName
	case "postgres":
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
	case "sqlite3":
		query = "PRAGMA index_list(" + tableName + ")"
	default:
		return []types.IndexInfo{}, fmt.Errorf("unsupported database type: %s", da.dbType)
	}

	rows, err := da.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer rows.Close()

	var indexes []types.IndexInfo
	for rows.Next() {
		var index types.IndexInfo

		switch da.dbType {
		case "mysql":
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
		case "postgres":
			var indexDef string
			err = rows.Scan(&index.Name, &indexDef, &index.IsUnique)
			if err != nil {
				continue
			}
			index.Type = "btree" // Default for PostgreSQL
			// Parse column names from indexdef if needed
			index.Columns = []string{"parsed_from_def"}
		case "sqlite3":
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

			// Get column names for this index
			colQuery := "PRAGMA index_info(" + index.Name + ")"
			colRows, colErr := da.db.Query(colQuery)
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

func (da *DatabaseAnalyzerImpl) GetConstraints(tableName string) ([]types.Constraint, error) {
	var query string
	var args []interface{}

	switch da.dbType {
	case "mysql":
		query = `
			SELECT 
				CONSTRAINT_NAME, CONSTRAINT_TYPE, COLUMN_NAME
			FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
			LEFT JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu 
				ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME 
				AND tc.TABLE_NAME = kcu.TABLE_NAME
			WHERE tc.TABLE_NAME = ? AND tc.TABLE_SCHEMA = DATABASE()`
		args = append(args, tableName)
	case "postgres":
		query = `
			SELECT 
				tc.constraint_name, tc.constraint_type, kcu.column_name
			FROM information_schema.table_constraints tc
			LEFT JOIN information_schema.key_column_usage kcu 
				ON tc.constraint_name = kcu.constraint_name 
				AND tc.table_name = kcu.table_name
			WHERE tc.table_name = $1`
		args = append(args, tableName)
	case "sqlite3":
		query = "PRAGMA foreign_key_list(" + tableName + ")"
	default:
		return []types.Constraint{}, fmt.Errorf("unsupported database type: %s", da.dbType)
	}

	rows, err := da.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query constraints: %w", err)
	}
	defer rows.Close()

	var constraints []types.Constraint
	for rows.Next() {
		var constraint types.Constraint

		switch da.dbType {
		case "mysql", "postgres":
			var columnName sql.NullString
			err = rows.Scan(&constraint.Name, &constraint.Type, &columnName)
			if err != nil {
				continue
			}
			if columnName.Valid {
				constraint.Columns = []string{columnName.String}
			}
		case "sqlite3":
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

func (da *DatabaseAnalyzerImpl) GetRelationships(tableName string) ([]types.Relationship, error) {
	var query string
	var args []interface{}

	switch da.dbType {
	case "mysql":
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
	case "postgres":
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
	case "sqlite3":
		query = "PRAGMA foreign_key_list(" + tableName + ")"
	default:
		return []types.Relationship{}, fmt.Errorf("unsupported database type: %s", da.dbType)
	}

	rows, err := da.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships: %w", err)
	}
	defer rows.Close()

	var relationships []types.Relationship
	for rows.Next() {
		var rel types.Relationship

		switch da.dbType {
		case "mysql", "postgres":
			var updateRule, deleteRule sql.NullString
			err = rows.Scan(&rel.SourceColumn, &rel.TargetTable,
				&rel.TargetColumn, &updateRule, &deleteRule)
			if err != nil {
				continue
			}
			rel.Type = "FOREIGN KEY"
		case "sqlite3":
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

func (da *DatabaseAnalyzerImpl) GetTableSize(tableName string) (string, error) {
	var query string
	var args []interface{}

	switch da.dbType {
	case "mysql":
		query = `
			SELECT 
				ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
			FROM information_schema.tables 
			WHERE table_schema = DATABASE() AND table_name = ?`
		args = append(args, tableName)
	case "postgres":
		query = `
			SELECT 
				pg_size_pretty(pg_total_relation_size($1)) AS size
			`
		args = append(args, tableName)
	case "sqlite3":
		query = `
			SELECT 
				COUNT(*) * 1024 as approx_bytes
			FROM pragma_table_info(?) 
			LIMIT 1`
		args = append(args, tableName)
	default:
		return "Unknown", fmt.Errorf("unsupported database type: %s", da.dbType)
	}

	var sizeResult interface{}
	err := da.db.QueryRow(query, args...).Scan(&sizeResult)
	if err != nil {
		log.Printf("Warning: Could not get table size for %s: %v", tableName, err)
		return "Unknown", nil
	}

	switch da.dbType {
	case "mysql":
		if size, ok := sizeResult.(float64); ok {
			return fmt.Sprintf("%.2f MB", size), nil
		}
	case "postgres":
		if size, ok := sizeResult.(string); ok {
			return size, nil
		}
	case "sqlite3":
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
