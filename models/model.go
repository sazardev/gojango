package models

import (
	"time"
)

// Model provides basic fields that all models should have (like Django's Model)
type Model struct {
	ID        uint      `json:"id" db:"id,primary_key,auto_increment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// BeforeCreate sets CreatedAt and UpdatedAt timestamps
func (m *Model) BeforeCreate() {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
}

// BeforeUpdate sets UpdatedAt timestamp
func (m *Model) BeforeUpdate() {
	m.UpdatedAt = time.Now()
}

// TableName returns the table name for the model (override in your models)
func (m *Model) TableName() string {
	return ""
}

// ModelInterface defines the interface that all models should implement
type ModelInterface interface {
	TableName() string
	BeforeCreate()
	BeforeUpdate()
}

// Field tags for database mapping
const (
	TagDB          = "db"
	TagJSON        = "json"
	TagPrimaryKey  = "primary_key"
	TagAutoIncr    = "auto_increment"
	TagNotNull     = "not_null"
	TagUnique      = "unique"
	TagIndex       = "index"
	TagDefault     = "default"
	TagSize        = "size"
	TagType        = "type"
)

// Common field types
const (
	TypeText      = "TEXT"
	TypeInteger   = "INTEGER"
	TypeReal      = "REAL"
	TypeBlob      = "BLOB"
	TypeDatetime  = "DATETIME"
	TypeBoolean   = "BOOLEAN"
	TypeVarchar   = "VARCHAR"
)

// ValidationError represents a model validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// Validator interface for model validation
type Validator interface {
	Validate() []ValidationError
}

// Example model structure that users can follow:
/*
type User struct {
	models.Model
	Name     string `json:"name" db:"name,not_null,size:100"`
	Email    string `json:"email" db:"email,unique,not_null,size:255"`
	Password string `json:"-" db:"password,not_null,size:255"`
	Active   bool   `json:"active" db:"active,default:true"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Validate() []models.ValidationError {
	var errors []models.ValidationError
	
	if u.Name == "" {
		errors = append(errors, models.ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	}
	
	if u.Email == "" {
		errors = append(errors, models.ValidationError{
			Field:   "email", 
			Message: "Email is required",
		})
	}
	
	return errors
}
*/
