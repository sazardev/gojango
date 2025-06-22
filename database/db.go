package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// MockDB is a simple in-memory database for testing
type MockDB struct {
	tables map[string][]map[string]interface{}
	nextID map[string]int
	mutex  sync.RWMutex
}

// DB wraps database connection with ORM-like functionality
type DB struct {
	Conn   *sql.DB // Exported for external access
	driver string
	mock   *MockDB // For testing without CGO
}

// Connect establishes database connection
func Connect(databaseURL string) (*DB, error) {
	// Simple URL parsing - in production you'd want more robust parsing
	var driver, dsn string

	if databaseURL == "" || strings.HasPrefix(databaseURL, "sqlite") {
		driver = "sqlite3"
		if databaseURL == "" {
			dsn = ":memory:"
		} else {
			dsn = strings.TrimPrefix(databaseURL, "sqlite://")
			if dsn == "" {
				dsn = ":memory:"
			}
		}
	} else {
		return nil, fmt.Errorf("unsupported database URL: %s", databaseURL)
	}

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DB{
		Conn:   conn,
		driver: driver,
	}, nil
}

// ConnectMock creates a mock database connection for testing
func ConnectMock() (*DB, error) {
	return &DB{
		Conn:   nil, // No real connection for mock
		driver: "mock",
		mock: &MockDB{
			tables: make(map[string][]map[string]interface{}),
			nextID: make(map[string]int),
		},
	}, nil
}

// AutoMigrate creates/updates table schema for the given model
func (db *DB) AutoMigrate(model interface{}) error {
	// Use mock database if available
	if db.mock != nil {
		return db.mock.AutoMigrate(model)
	}

	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)

	// Handle pointer types
	if modelType.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
		modelType = modelType.Elem()
	}

	// Get table name
	tableName := db.getTableName(model)
	if tableName == "" {
		tableName = strings.ToLower(modelType.Name()) + "s"
	}

	// Build CREATE TABLE statement
	var columns []string

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		columnDef := db.buildColumnDefinition(field, dbTag)
		if columnDef != "" {
			columns = append(columns, columnDef)
		}
	}

	if len(columns) == 0 {
		return fmt.Errorf("no database columns found for model %T", model)
	}

	createSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n)",
		tableName, strings.Join(columns, ",\n  "))

	_, err := db.Conn.Exec(createSQL)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %v", tableName, err)
	}

	return nil
}

// buildColumnDefinition creates column definition from field and tag
func (db *DB) buildColumnDefinition(field reflect.StructField, dbTag string) string {
	parts := strings.Split(dbTag, ",")
	columnName := parts[0]

	if columnName == "" {
		return ""
	}

	// Determine column type based on Go type
	var columnType string
	switch field.Type.Kind() {
	case reflect.String:
		columnType = "TEXT"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		columnType = "INTEGER"
	case reflect.Float32, reflect.Float64:
		columnType = "REAL"
	case reflect.Bool:
		columnType = "BOOLEAN"
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.Uint8 {
			columnType = "BLOB"
		} else {
			columnType = "TEXT"
		}
	default:
		if field.Type == reflect.TypeOf(time.Time{}) {
			columnType = "DATETIME"
		} else {
			columnType = "TEXT"
		}
	}

	// Parse additional options
	var constraints []string

	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		switch {
		case part == "primary_key":
			constraints = append(constraints, "PRIMARY KEY")
		case part == "auto_increment":
			constraints = append(constraints, "AUTOINCREMENT")
		case part == "not_null":
			constraints = append(constraints, "NOT NULL")
		case part == "unique":
			constraints = append(constraints, "UNIQUE")
		case strings.HasPrefix(part, "default:"):
			defaultVal := strings.TrimPrefix(part, "default:")
			constraints = append(constraints, "DEFAULT "+defaultVal)
		case strings.HasPrefix(part, "size:"):
			size := strings.TrimPrefix(part, "size:")
			if columnType == "TEXT" {
				columnType = fmt.Sprintf("VARCHAR(%s)", size)
			}
		case strings.HasPrefix(part, "type:"):
			columnType = strings.TrimPrefix(part, "type:")
		}
	}

	definition := fmt.Sprintf("%s %s", columnName, columnType)
	if len(constraints) > 0 {
		definition += " " + strings.Join(constraints, " ")
	}

	return definition
}

// GetTableName extracts table name from model (exported for external use)
func (db *DB) GetTableName(model interface{}) string {
	return db.getTableName(model)
}

// getTableName extracts table name from model
func (db *DB) getTableName(model interface{}) string {
	if tableNamer, ok := model.(interface{ TableName() string }); ok {
		return tableNamer.TableName()
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	return strings.ToLower(modelType.Name()) + "s"
}

// Create inserts a new record
func (db *DB) Create(model interface{}) error {
	// Use mock database if available
	if db.mock != nil {
		return db.mock.Create(model)
	}

	// Call BeforeCreate hook if available
	if beforeCreator, ok := model.(interface{ BeforeCreate() }); ok {
		beforeCreator.BeforeCreate()
	}

	tableName := db.getTableName(model)

	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
		modelType = modelType.Elem()
	}

	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		if !field.IsExported() {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		columnName := strings.Split(dbTag, ",")[0]

		// Skip auto-increment primary keys
		if strings.Contains(dbTag, "auto_increment") {
			continue
		}

		columns = append(columns, columnName)
		placeholders = append(placeholders, "?")
		values = append(values, fieldValue.Interface())
	}

	if len(columns) == 0 {
		return fmt.Errorf("no columns to insert for model %T", model)
	}

	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	result, err := db.Conn.Exec(insertSQL, values...)
	if err != nil {
		return fmt.Errorf("failed to insert record: %v", err)
	}

	// Set the ID if it's an auto-increment field
	if lastID, err := result.LastInsertId(); err == nil && lastID > 0 {
		db.setIDField(model, lastID)
	}

	return nil
}

// FindAll retrieves all records of a model type
func (db *DB) FindAll(model interface{}) (interface{}, error) {
	// Use mock database if available
	if db.mock != nil {
		return db.mock.FindAll(model)
	}

	tableName := db.getTableName(model)

	selectSQL := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Conn.Query(selectSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %v", err)
	}
	defer rows.Close()

	return db.scanRows(rows, model)
}

// FindByID finds a record by ID
func (db *DB) FindByID(model interface{}, id string) error {
	// Use mock database if available
	if db.mock != nil {
		return db.mock.FindByID(model, id)
	}

	tableName := db.getTableName(model)

	selectSQL := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
	row := db.Conn.QueryRow(selectSQL, id)

	return db.scanRow(row, model)
}

// Update updates a record by ID
func (db *DB) Update(model interface{}, id string) error {
	// Call BeforeUpdate hook if available
	if beforeUpdater, ok := model.(interface{ BeforeUpdate() }); ok {
		beforeUpdater.BeforeUpdate()
	}

	tableName := db.getTableName(model)

	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
		modelType = modelType.Elem()
	}

	var setParts []string
	var values []interface{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		if !field.IsExported() {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		columnName := strings.Split(dbTag, ",")[0]

		// Skip primary key and auto-increment fields
		if strings.Contains(dbTag, "primary_key") || strings.Contains(dbTag, "auto_increment") {
			continue
		}

		setParts = append(setParts, columnName+" = ?")
		values = append(values, fieldValue.Interface())
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no columns to update for model %T", model)
	}

	values = append(values, id)
	updateSQL := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		tableName, strings.Join(setParts, ", "))

	_, err := db.Conn.Exec(updateSQL, values...)
	if err != nil {
		return fmt.Errorf("failed to update record: %v", err)
	}

	return nil
}

// Delete deletes a record by ID
func (db *DB) Delete(model interface{}, id string) error {
	tableName := db.getTableName(model)

	deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := db.Conn.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}

// setIDField sets the ID field of a model (helper for auto-increment)
func (db *DB) setIDField(model interface{}, id int64) {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	if !modelValue.CanSet() {
		return
	}

	// Look for ID field
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Type().Field(i)
		fieldValue := modelValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		dbTag := field.Tag.Get("db")
		if strings.Contains(dbTag, "primary_key") && strings.Contains(dbTag, "auto_increment") {
			switch fieldValue.Kind() {
			case reflect.Uint, reflect.Uint32, reflect.Uint64:
				fieldValue.SetUint(uint64(id))
			case reflect.Int, reflect.Int32, reflect.Int64:
				fieldValue.SetInt(id)
			}
			break
		}
	}
}

// ScanRows scans multiple rows into a slice of models (exported for external use)
func (db *DB) ScanRows(rows *sql.Rows, model interface{}) (interface{}, error) {
	return db.scanRows(rows, model)
}

// scanRows scans multiple rows into a slice of models
func (db *DB) scanRows(rows *sql.Rows, model interface{}) (interface{}, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	sliceType := reflect.SliceOf(reflect.PtrTo(modelType))
	results := reflect.MakeSlice(sliceType, 0, 0)

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		newModel := reflect.New(modelType)

		if err := db.scanRowIntoModel(rows, columns, newModel.Interface()); err != nil {
			return nil, err
		}

		results = reflect.Append(results, newModel)
	}

	return results.Interface(), nil
}

// scanRow scans a single row into a model
func (db *DB) scanRow(row *sql.Row, model interface{}) error {
	// For single row, we need to get columns differently
	// This is a simplified implementation
	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
		modelType = modelType.Elem()
	}

	var scanValues []interface{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		if !field.IsExported() {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		scanValues = append(scanValues, fieldValue.Addr().Interface())
	}

	return row.Scan(scanValues...)
}

// scanRowIntoModel scans a row into a model with column mapping
func (db *DB) scanRowIntoModel(rows *sql.Rows, columns []string, model interface{}) error {
	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
		modelType = modelType.Elem()
	}

	// Create a map of column names to field indices
	columnMap := make(map[string]int)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			columnName := strings.Split(dbTag, ",")[0]
			columnMap[columnName] = i
		}
	}

	// Prepare scan destinations
	scanDests := make([]interface{}, len(columns))
	for i, column := range columns {
		if fieldIndex, exists := columnMap[column]; exists {
			fieldValue := modelValue.Field(fieldIndex)
			scanDests[i] = fieldValue.Addr().Interface()
		} else {
			// Use a discard variable for unknown columns
			var discard interface{}
			scanDests[i] = &discard
		}
	}

	return rows.Scan(scanDests...)
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.Conn.Close()
}

// MockAutoMigrate simulates table creation
func (mdb *MockDB) AutoMigrate(model interface{}) error {
	tableName := mdb.getTableName(model)

	mdb.mutex.Lock()
	defer mdb.mutex.Unlock()

	if _, exists := mdb.tables[tableName]; !exists {
		mdb.tables[tableName] = make([]map[string]interface{}, 0)
		mdb.nextID[tableName] = 1
	}

	return nil
}

// MockCreate simulates record creation
func (mdb *MockDB) Create(model interface{}) error {
	tableName := mdb.getTableName(model)

	mdb.mutex.Lock()
	defer mdb.mutex.Unlock()

	// Initialize table if it doesn't exist
	if _, exists := mdb.tables[tableName]; !exists {
		mdb.AutoMigrate(model)
	}

	// Convert model to map
	record := mdb.modelToMap(model)

	// Set ID if not set
	if _, hasID := record["id"]; !hasID {
		record["id"] = mdb.nextID[tableName]
		mdb.nextID[tableName]++

		// Set ID back to model if possible
		if v := reflect.ValueOf(model); v.Kind() == reflect.Ptr {
			elem := v.Elem()
			if idField := elem.FieldByName("ID"); idField.IsValid() && idField.CanSet() {
				idField.SetUint(uint64(record["id"].(int)))
			}
		}
	}

	mdb.tables[tableName] = append(mdb.tables[tableName], record)
	return nil
}

// MockFindAll simulates finding all records
func (mdb *MockDB) FindAll(model interface{}) (interface{}, error) {
	tableName := mdb.getTableName(model)

	mdb.mutex.RLock()
	defer mdb.mutex.RUnlock()

	records, exists := mdb.tables[tableName]
	if !exists {
		// Return empty slice of correct type
		sliceType := reflect.SliceOf(reflect.PtrTo(reflect.TypeOf(model).Elem()))
		return reflect.MakeSlice(sliceType, 0, 0).Interface(), nil
	}

	// Create slice of pointers to models
	modelType := reflect.TypeOf(model).Elem()
	sliceType := reflect.SliceOf(reflect.PtrTo(modelType))
	results := reflect.MakeSlice(sliceType, len(records), len(records))

	for i, record := range records {
		// Create new model instance
		newModel := reflect.New(modelType)

		// Fill model with data
		mdb.mapToModel(record, newModel.Interface())

		results.Index(i).Set(newModel)
	}

	return results.Interface(), nil
}

// MockFindByID simulates finding a record by ID
func (mdb *MockDB) FindByID(model interface{}, id string) error {
	tableName := mdb.getTableName(model)

	mdb.mutex.RLock()
	defer mdb.mutex.RUnlock()

	records, exists := mdb.tables[tableName]
	if !exists {
		return fmt.Errorf("record not found")
	}

	// Convert ID to int for comparison
	var targetID int
	if _, err := fmt.Sscanf(id, "%d", &targetID); err != nil {
		return fmt.Errorf("invalid ID format")
	}

	for _, record := range records {
		if recordID, hasID := record["id"]; hasID && recordID == targetID {
			return mdb.mapToModel(record, model)
		}
	}

	return fmt.Errorf("record not found")
}

func (mdb *MockDB) mapToModel(data map[string]interface{}, model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("model must be a pointer")
	}

	elem := v.Elem()
	t := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := t.Field(i)
		value := elem.Field(i)

		if !value.CanSet() {
			continue
		}

		fieldName := strings.ToLower(field.Name)

		if dataValue, exists := data[fieldName]; exists && dataValue != nil {
			dataVal := reflect.ValueOf(dataValue)
			if dataVal.Type().ConvertibleTo(value.Type()) {
				value.Set(dataVal.Convert(value.Type()))
			}
		}
	}

	return nil
}

// Helper methods for MockDB
func (mdb *MockDB) getTableName(model interface{}) string {
	if tabler, ok := model.(interface{ TableName() string }); ok {
		return tabler.TableName()
	}

	// Default to struct name
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name()) + "s"
}

func (mdb *MockDB) modelToMap(model interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.CanInterface() {
			fieldName := strings.ToLower(field.Name)
			result[fieldName] = value.Interface()
		}
	}

	return result
}

// Query executes a query with mock support
func (db *DB) Query(query string, args ...interface{}) (interface{}, error) {
	if db.mock != nil {
		// For now, return empty results for mock queries
		// In a full implementation, you'd parse the SQL and simulate the query
		return []interface{}{}, nil
	}

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanRows(rows, nil)
}
