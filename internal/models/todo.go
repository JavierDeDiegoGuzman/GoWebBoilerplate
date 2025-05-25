package models

import (
    "database/sql"
    "fmt"
    "time"
    "GoWebBoilerplate/internal/utils"
)

// Todo representa un item en la lista de tareas
type Todo struct {
    ID        string    `db:"id,primary_key"`
    Task      string    `db:"task,not_null"`
    Done      bool      `db:"done"`
    CreatedAt string    `db:"created_at"`
    UpdatedAt string    `db:"updated_at"`
}

func init() {
    // Registrar el modelo Todo automáticamente
    utils.RegisterModel("todos", Todo{})
}

// NewTodo crea una nueva instancia de Todo
func NewTodo(task string) *Todo {
    now := time.Now().Format(time.RFC3339)
    return &Todo{
        Task:      task,
        Done:      false,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

// Save guarda el todo en la base de datos
func (t *Todo) Save() error {
    if t.ID == "" {
        // Es un nuevo todo, hacer INSERT
        return t.insert()
    }
    // Es un todo existente, hacer UPDATE
    return t.update()
}

// insert crea un nuevo todo en la base de datos
func (t *Todo) insert() error {
    // Generar un ID único
    id, err := utils.GenerateUniqueID("todos")
    if err != nil {
        return fmt.Errorf("error generating unique ID: %w", err)
    }
    t.ID = id

    query := `
        INSERT INTO todos (id, task, done, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `
    
    now := time.Now().Format(time.RFC3339)
    t.CreatedAt = now
    t.UpdatedAt = now

    _, err = utils.GetDB().Exec(query, 
        t.ID,
        t.Task, 
        t.Done, 
        t.CreatedAt, 
        t.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("error creating todo: %w", err)
    }

    return nil
}

// update actualiza un todo existente
func (t *Todo) update() error {
    query := `
        UPDATE todos 
        SET task = ?, done = ?, updated_at = ?
        WHERE id = ?
    `
    
    t.UpdatedAt = time.Now().Format(time.RFC3339)
    
    result, err := utils.GetDB().Exec(query,
        t.Task,
        t.Done,
        t.UpdatedAt,
        t.ID,
    )
    if err != nil {
        return fmt.Errorf("error updating todo: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }
    if rows == 0 {
        return sql.ErrNoRows
    }

    return nil
}

// GetByID obtiene un todo por su ID
func (t *Todo) GetByID(id string) error {
    return utils.GetDB().Get(t, 
        "SELECT * FROM todos WHERE id = ?", 
        id,
    )
}

// GetAll obtiene todos los todos
func (t *Todo) GetAll() ([]Todo, error) {
    var todos []Todo
    err := utils.GetDB().Select(&todos, 
        "SELECT * FROM todos ORDER BY created_at DESC",
    )
    return todos, err
}

// Delete elimina un todo
func (t *Todo) Delete() error {
    result, err := utils.GetDB().Exec(
        "DELETE FROM todos WHERE id = ?", 
        t.ID,
    )
    if err != nil {
        return fmt.Errorf("error deleting todo: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }
    if rows == 0 {
        return sql.ErrNoRows
    }

    return nil
}
