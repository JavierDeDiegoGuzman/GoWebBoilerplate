package utils

import (
    "context"
    "fmt"
    "log"
    "reflect"
    "strings"
    "time"

    _ "github.com/glebarez/sqlite"
    "github.com/jmoiron/sqlx"
)

// dbManager maneja la conexión a la base de datos
type dbManager struct {
    db *sqlx.DB
}

var (
    manager *dbManager
    defaultCtx = context.Background()
    ModelTableMapping = make(map[string]interface{})
)

// InitDB inicializa la conexión a la base de datos
func InitDB(dsn string) error {
    db, err := sqlx.Connect("sqlite", dsn)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Configurar la conexión
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    manager = &dbManager{db: db}
    return nil
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() error {
    if manager != nil {
        return manager.db.Close()
    }
    return nil
}

// GetDB retorna la conexión a la base de datos
func GetDB() *sqlx.DB {
    if manager == nil {
        log.Fatal("Database not initialized")
    }
    return manager.db
}

// WithTransaction ejecuta una función dentro de una transacción
func WithTransaction(fn func(*sqlx.Tx) error) error {
    if manager == nil {
        return fmt.Errorf("database not initialized")
    }

    tx, err := manager.db.Beginx()
    if err != nil {
        return err
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()

    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

// RegisterModel registra un modelo en la base de datos
func RegisterModel(tableName string, model interface{}) {
    ModelTableMapping[tableName] = model
}

// getSQLiteType convierte un tipo de Go a un tipo SQLite
func getSQLiteType(t reflect.Type) string {
    switch t.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return "INTEGER"
    case reflect.Float32, reflect.Float64:
        return "REAL"
    case reflect.Bool:
        return "BOOLEAN"
    case reflect.String:
        return "TEXT"
    case reflect.Slice:
        if t.Elem().Kind() == reflect.Uint8 { // []byte
            return "BLOB"
        }
    }
    return "TEXT" // default
}

// GenerateTableSchema genera el esquema CREATE TABLE para un struct
func GenerateTableSchema(structType interface{}, tableName string) string {
    t := reflect.TypeOf(structType)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }

    var columns []string
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        if field.Anonymous {
            continue
        }

        dbTag := field.Tag.Get("db")
        if dbTag == "" || dbTag == "-" {
            continue
        }

        parts := strings.Split(dbTag, ",")
        columnName := parts[0]

        sqlType := getSQLiteType(field.Type)
        column := fmt.Sprintf("%s %s", columnName, sqlType)

        // Manejar tags especiales
        if strings.Contains(dbTag, "primary_key") {
            column += " PRIMARY KEY"
            if sqlType == "INTEGER" {
                column += " AUTOINCREMENT"
            }
        }
        if strings.Contains(dbTag, "not_null") {
            column += " NOT NULL"
        }
        if strings.Contains(dbTag, "unique") {
            column += " UNIQUE"
        }

        columns = append(columns, column)
    }

    if len(columns) == 0 {
        log.Printf("Warning: No columns found for table %s", tableName)
        return ""
    }

    return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        %s
    );`, tableName, strings.Join(columns, ",\n\t\t"))
}

// CreateTables crea todas las tablas necesarias si no existen
func CreateTables() error {
    for tableName, model := range ModelTableMapping {
        log.Printf("Creating table: %s", tableName)
        schema := GenerateTableSchema(model, tableName)
        if schema == "" {
            return fmt.Errorf("failed to generate schema for table %s: no columns found", tableName)
        }

        log.Printf("Creating table with schema:\n%s", schema)
        if err := CreateTableFromStruct(model, tableName); err != nil {
            return fmt.Errorf("failed to create table %s: %w", tableName, err)
        }
    }
    return nil
}

// CreateTableFromStruct crea una tabla basada en un struct
func CreateTableFromStruct(structType interface{}, tableName string) error {
    schema := GenerateTableSchema(structType, tableName)
    if schema == "" {
        return fmt.Errorf("empty schema generated for table %s", tableName)
    }

    _, err := manager.db.Exec(schema)
    if err != nil {
        return fmt.Errorf("failed to create table %s: %w", tableName, err)
    }
    return nil
}

// GenerateUniqueID genera un ID único para una tabla usando randomblob
func GenerateUniqueID(tableName string) (string, error) {
    // Intentar generar un ID único hasta 5 veces
    for i := 0; i < 5; i++ {
        // Generar un ID aleatorio de 16 bytes (32 caracteres en hex)
        var id string
        err := manager.db.Get(&id, "SELECT lower(hex(randomblob(16)))")
        if err != nil {
            return "", fmt.Errorf("error generating random ID: %w", err)
        }

        // Verificar si el ID ya existe
        var exists bool
        err = manager.db.Get(&exists, fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = ?)", tableName), id)
        if err != nil {
            return "", fmt.Errorf("error checking ID existence: %w", err)
        }

        if !exists {
            return id, nil
        }
    }

    return "", fmt.Errorf("failed to generate unique ID after 5 attempts")
} 
