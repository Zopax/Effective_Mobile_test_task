package service

import (
	"context"
	"time"

	"effective-mobile-task/internal/client"
	"effective-mobile-task/internal/models"
	"effective-mobile-task/internal/repository"
	"github.com/google/uuid"
)

type PersonService struct {
	repo     *repository.PersonRepository
	enricher *client.EnrichmentClient
}

func NewPersonService(r *repository.PersonRepository, e *client.EnrichmentClient) *PersonService {
	return &PersonService{repo: r, enricher: e}
}

func (s *PersonService) Create(ctx context.Context, req models.CreatePersonRequest) (*models.Person, error) {
	age, err := s.enricher.GetAge(req.Name)
	if err != nil {
		return nil, err
	}
	gender, err := s.enricher.GetGender(req.Name)
	if err != nil {
		return nil, err
	}
	nationality, err := s.enricher.GetNationality(req.Name)
	if err != nil {
		return nil, err
	}

	person := &models.Person{
		ID:          uuid.New(),
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.repo.Create(ctx, *person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (s *PersonService) Update(ctx context.Context, id uuid.UUID, req models.UpdatePersonRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *PersonService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *PersonService) GetByID(ctx context.Context, id uuid.UUID) (*models.Person, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PersonService) List(ctx context.Context, limit, offset int, filters map[string]string) ([]models.Person, error) {
	return s.repo.GetAll(ctx, limit, offset, filters)
}
