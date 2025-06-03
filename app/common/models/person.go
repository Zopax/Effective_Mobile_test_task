package models

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Surname     string    `json:"surname" db:"surname"`
	Patronymic  *string   `json:"patronymic,omitempty" db:"patronymic"`
	Age         int       `json:"age" db:"age"`
	Gender      string    `json:"gender" db:"gender"`
	Nationality string    `json:"nationality" db:"nationality"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
