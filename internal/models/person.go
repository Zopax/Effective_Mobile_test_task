package models

import (
	"time"

	"github.com/google/uuid"
)

// Person
type Person struct {
	ID          uuid.UUID `json:"id" example:"1e8c72e6-3c77-4b9b-b44d-1b0e44c3c0b9"`
	Name        string    `json:"name" example:"Dmitriy"`
	Surname     string    `json:"surname" example:"Ushakov"`
	Patronymic  string    `json:"patronymic,omitempty" example:"Vasilevich"`
	Age         int       `json:"age" example:"30"`
	Gender      string    `json:"gender" example:"male"`
	Nationality string    `json:"nationality" example:"RU"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePersonRequest
type CreatePersonRequest struct {
	Name       string `json:"name" example:"Dmitriy"`
	Surname    string `json:"surname" example:"Ushakov"`
	Patronymic string `json:"patronymic,omitempty" example:"Vasilevich"`
}

// UpdatePersonRequest
type UpdatePersonRequest struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
}
