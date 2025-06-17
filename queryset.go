package gojango

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	
	"gojango/database"
)

// QuerySet provides Django-like query capabilities
type QuerySet struct {
	db        *database.DB
	model     interface{}
	modelType reflect.Type
	tableName string
	where     []string
	args      []interface{}
	orderBy   string
	limit     int
	offset    int
}

// NewQuerySet creates a new QuerySet for a model
func (app *App) NewQuerySet(model interface{}) *QuerySet {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	
	qs := &QuerySet{
		db:        app.db,
		model:     model,
		modelType: modelType,
		tableName: app.db.GetTableName(model),
	}
	
	return qs
}

// Filter adds WHERE conditions (Django-like)
func (qs *QuerySet) Filter(field string, value interface{}) *QuerySet {
	// Create a copy to avoid mutating the original
	newQS := *qs
	newQS.where = make([]string, len(qs.where))
	copy(newQS.where, qs.where)
	newQS.args = make([]interface{}, len(qs.args))
	copy(newQS.args, qs.args)
	
	// Parse Django-style field lookups
	parts := strings.Split(field, "__")
	fieldName := parts[0]
	lookup := "exact"
	
	if len(parts) > 1 {
		lookup = parts[1]
	}
	
	var condition string
	switch lookup {
	case "exact":
		condition = fieldName + " = ?"
	case "iexact":
		condition = "LOWER(" + fieldName + ") = LOWER(?)"
	case "contains":
		condition = fieldName + " LIKE ?"
		value = "%" + fmt.Sprintf("%v", value) + "%"
	case "icontains":
		condition = "LOWER(" + fieldName + ") LIKE LOWER(?)"
		value = "%" + fmt.Sprintf("%v", value) + "%"
	case "startswith":
		condition = fieldName + " LIKE ?"
		value = fmt.Sprintf("%v", value) + "%"
	case "endswith":
		condition = fieldName + " LIKE ?"
		value = "%" + fmt.Sprintf("%v", value)
	case "gt":
		condition = fieldName + " > ?"
	case "gte":
		condition = fieldName + " >= ?"
	case "lt":
		condition = fieldName + " < ?"
	case "lte":
		condition = fieldName + " <= ?"
	case "in":
		// Handle IN queries
		if slice := reflect.ValueOf(value); slice.Kind() == reflect.Slice {
			placeholders := make([]string, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				placeholders[i] = "?"
				newQS.args = append(newQS.args, slice.Index(i).Interface())
			}
			condition = fieldName + " IN (" + strings.Join(placeholders, ",") + ")"
			// Don't add value to args since we already added individual items
			goto skipValueAdd
		}
		condition = fieldName + " = ?"
	case "isnull":
		if value.(bool) {
			condition = fieldName + " IS NULL"
		} else {
			condition = fieldName + " IS NOT NULL"
		}
		// Don't add value to args for NULL checks
		goto skipValueAdd
	default:
		condition = fieldName + " = ?"
	}
	
	newQS.where = append(newQS.where, condition)
	newQS.args = append(newQS.args, value)
	
skipValueAdd:
	return &newQS
}

// Exclude adds WHERE NOT conditions
func (qs *QuerySet) Exclude(field string, value interface{}) *QuerySet {
	// Similar to Filter but with NOT
	newQS := qs.Filter(field, value)
	// Modify the last condition to be NOT
	if len(newQS.where) > 0 {
		lastIndex := len(newQS.where) - 1
		newQS.where[lastIndex] = "NOT (" + newQS.where[lastIndex] + ")"
	}
	return newQS
}

// OrderBy adds ORDER BY clause
func (qs *QuerySet) OrderBy(field string) *QuerySet {
	newQS := *qs
	
	// Handle Django-style ordering
	if strings.HasPrefix(field, "-") {
		newQS.orderBy = strings.TrimPrefix(field, "-") + " DESC"
	} else {
		newQS.orderBy = field + " ASC"
	}
	
	return &newQS
}

// Limit adds LIMIT clause
func (qs *QuerySet) Limit(limit int) *QuerySet {
	newQS := *qs
	newQS.limit = limit
	return &newQS
}

// Offset adds OFFSET clause
func (qs *QuerySet) Offset(offset int) *QuerySet {
	newQS := *qs
	newQS.offset = offset
	return &newQS
}

// All executes the query and returns all results
func (qs *QuerySet) All() (interface{}, error) {
	sql := qs.buildSQL()
	
	rows, err := qs.db.conn.Query(sql, qs.args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()
	
	return qs.db.scanRows(rows, qs.model)
}

// First returns the first result
func (qs *QuerySet) First() (interface{}, error) {
	limitedQS := qs.Limit(1)
	results, err := limitedQS.All()
	if err != nil {
		return nil, err
	}
	
	// Extract first item from slice
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() == reflect.Slice && resultsValue.Len() > 0 {
		return resultsValue.Index(0).Interface(), nil
	}
	
	return nil, fmt.Errorf("no results found")
}

// Count returns the count of matching records
func (qs *QuerySet) Count() (int, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) FROM %s", qs.tableName)
	
	if len(qs.where) > 0 {
		sql += " WHERE " + strings.Join(qs.where, " AND ")
	}
	
	var count int
	err := qs.db.conn.QueryRow(sql, qs.args...).Scan(&count)
	return count, err
}

// Exists checks if any records match the query
func (qs *QuerySet) Exists() (bool, error) {
	count, err := qs.Count()
	return count > 0, err
}

// buildSQL builds the complete SQL query
func (qs *QuerySet) buildSQL() string {
	sql := fmt.Sprintf("SELECT * FROM %s", qs.tableName)
	
	if len(qs.where) > 0 {
		sql += " WHERE " + strings.Join(qs.where, " AND ")
	}
	
	if qs.orderBy != "" {
		sql += " ORDER BY " + qs.orderBy
	}
	
	if qs.limit > 0 {
		sql += " LIMIT " + strconv.Itoa(qs.limit)
	}
	
	if qs.offset > 0 {
		sql += " OFFSET " + strconv.Itoa(qs.offset)
	}
	
	return sql
}

// Update updates matching records
func (qs *QuerySet) Update(data map[string]interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to update")
	}
	
	var setParts []string
	var args []interface{}
	
	for field, value := range data {
		setParts = append(setParts, field+" = ?")
		args = append(args, value)
	}
	
	sql := fmt.Sprintf("UPDATE %s SET %s", qs.tableName, strings.Join(setParts, ", "))
	
	if len(qs.where) > 0 {
		sql += " WHERE " + strings.Join(qs.where, " AND ")
		args = append(args, qs.args...)
	}
	
	_, err := qs.db.conn.Exec(sql, args...)
	return err
}

// Delete deletes matching records
func (qs *QuerySet) Delete() error {
	sql := fmt.Sprintf("DELETE FROM %s", qs.tableName)
	
	if len(qs.where) > 0 {
		sql += " WHERE " + strings.Join(qs.where, " AND ")
	}
	
	_, err := qs.db.conn.Exec(sql, qs.args...)
	return err
}

// ToJSON converts results to JSON
func (qs *QuerySet) ToJSON() (string, error) {
	results, err := qs.All()
	if err != nil {
		return "", err
	}
	
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	
	return string(jsonBytes), nil
}
