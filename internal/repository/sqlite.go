// Package repository
package repository

import (
	"database/sql"

	"flowtask/internal/models"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) GetAll(energyFilter string) ([]models.Task, error) {
	query := "SELECT id, title, description, energy, category, completed, created_at FROM tasks"

	var rows *sql.Rows
	var err error

	if energyFilter != "" && energyFilter != "all" {
		query += " WHERE energy = ?"
		rows, err = r.db.Query(query, energyFilter)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Energy, &t.Category, &t.Completed, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *SQLiteRepository) Create(t *models.Task) error {
	query := `INSERT INTO tasks (title, description, energy, category, completed) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, t.Title, t.Description, t.Energy, t.Category, t.Completed)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (r *SQLiteRepository) ToggleComplete(id int64) error {
	_, err := r.db.Exec("UPDATE tasks SET completed = CASE WHEN completed = 1 THEN 0 ELSE 1 END WHERE id = ?", id)
	return err
}

// Получить все категории с учетом порядка
func (r *SQLiteRepository) GetCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, position FROM categories ORDER BY position ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Position); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// Создать категорию (ставим в самый конец)
func (r *SQLiteRepository) CreateCategory(name string) (*models.Category, error) {
	var count int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count); err != nil {
		return nil, err
	}

	if count >= 10 {
		return nil, sql.ErrNoRows
	}

	// Узнаем максимальную позицию
	var maxPos int
	r.db.QueryRow("SELECT COALESCE(MAX(position), 0) FROM categories").Scan(&maxPos)

	result, err := r.db.Exec("INSERT INTO categories (name, position) VALUES (?, ?)", name, maxPos+1)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Category{ID: id, Name: name, Position: int64(maxPos + 1)}, nil
}

// Обновление порядка категорий
func (r *SQLiteRepository) UpdateCategoriesOrder(ids []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE categories SET position = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for idx, id := range ids {
		if _, err := stmt.Exec(idx+1, id); err != nil {
			return err
		}
	}

	return tx.Commit()
}
