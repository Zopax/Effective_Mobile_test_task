package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"effective-mobile-task/internal/models"
	"github.com/google/uuid"
)

type PersonRepository struct {
	db *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(ctx context.Context, p models.Person) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO persons (id, name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		p.ID, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Nationality)
	return err
}

func (r *PersonRepository) Update(ctx context.Context, id uuid.UUID, update models.UpdatePersonRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argPos := 1

	if update.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *update.Name)
		argPos++
	}
	if update.Surname != nil {
		setParts = append(setParts, fmt.Sprintf("surname = $%d", argPos))
		args = append(args, *update.Surname)
		argPos++
	}
	if update.Patronymic != nil {
		setParts = append(setParts, fmt.Sprintf("patronymic = $%d", argPos))
		args = append(args, *update.Patronymic)
		argPos++
	}

	if len(setParts) == 0 {
		return nil
	}

	setParts = append(setParts, "updated_at = NOW()")

	query := fmt.Sprintf("UPDATE persons SET %s WHERE id = $%d", strings.Join(setParts, ", "), argPos)
	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PersonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM persons WHERE id = $1", id)
	return err
}

func (r *PersonRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Person, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at
		FROM persons WHERE id = $1`, id)

	var p models.Person
	err := row.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PersonRepository) GetAll(ctx context.Context, limit, offset int, filters map[string]string) ([]models.Person, error) {
	query := `SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at
			  FROM persons WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if name, ok := filters["name"]; ok {
		query += fmt.Sprintf(" AND name ILIKE $%d", argID)
		args = append(args, "%"+name+"%")
		argID++
	}
	if surname, ok := filters["surname"]; ok {
		query += fmt.Sprintf(" AND surname ILIKE $%d", argID)
		args = append(args, "%"+surname+"%")
		argID++
	}
	if gender, ok := filters["gender"]; ok {
		query += fmt.Sprintf(" AND gender = $%d", argID)
		args = append(args, gender)
		argID++
	}
	if nationality, ok := filters["nationality"]; ok {
		query += fmt.Sprintf(" AND nationality = $%d", argID)
		args = append(args, nationality)
		argID++
	}
	if age, ok := filters["age"]; ok {
		query += fmt.Sprintf(" AND age = $%d", argID)
		args = append(args, age)
		argID++
	}
	//фильтрация по возрасту возраст больше чем
	if ageGT, ok := filters["age_gt"]; ok {
		query += fmt.Sprintf(" AND age > $%d", argID)
		args = append(args, ageGT)
		argID++
	}
	//фильтрация по возрасту возраст меньше чем
	if ageLT, ok := filters["age_lt"]; ok {
		query += fmt.Sprintf(" AND age < $%d", argID)
		args = append(args, ageLT)
		argID++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		people = append(people, p)
	}

	return people, nil
}
