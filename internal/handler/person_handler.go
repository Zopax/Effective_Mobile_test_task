package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"effective-mobile-task/internal/models"
	"effective-mobile-task/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PersonHandler struct {
	service *service.PersonService
	logger  *zap.Logger
}

func NewPersonHandler(s *service.PersonService, logger *zap.Logger) *PersonHandler {
	return &PersonHandler{service: s, logger: logger}
}

// Create godoc
// @Summary Создание нового человека
// @Description Обогащает ФИО через внешние API и сохраняет в БД
// @Tags persons
// @Accept json
// @Produce json
// @Param input body models.CreatePersonRequest true "Данные человека"
// @Success 201 {object} models.Person
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /persons [post]
func (h *PersonHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode create request", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	person, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create person", zap.Error(err))
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Created person", zap.String("id", person.ID.String()))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

// List godoc
// @Summary Получить список людей
// @Description Получить список людей с фильтрами и пагинацией
// @Tags persons
// @Produce json
// @Param limit query int false "Лимит"
// @Param offset query int false "Смещение"
// @Param name query string false "Имя"
// @Param surname query string false "Фамилия"
// @Param age query int false "Возраст"
// @Param age_gt query int false "Возраст больше"
// @Param age_lt query int false "Возраст меньше"
// @Param gender query string false "Пол"
// @Param nationality query string false "Национальность"
// @Success 200 {array} models.Person
// @Failure 500 {object} string
// @Router /persons [get]
func (h *PersonHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	filters := map[string]string{}
	for _, key := range []string{"name", "surname", "gender", "nationality", "age", "age_gt", "age_lt"} {
		if val := q.Get(key); val != "" {
			filters[key] = val
		}
	}

	people, err := h.service.List(r.Context(), limit, offset, filters)
	if err != nil {
		h.logger.Error("Failed to list persons", zap.Error(err))
		http.Error(w, "Failed to list persons", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

// GetById godoc
// @Summary Получить одного человека по ID
// @Tags persons
// @Produce json
// @Param id path string true "UUID человека"
// @Success 200 {object} models.Person
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /persons/{id} [get]

func (h *PersonHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid UUID", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	person, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get person by ID", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// Update godoc
// @Summary Обновить данные человека по ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "UUID человека"
// @Param input body models.CreatePersonRequest true "Новые данные"
// @Success 204 {object} models.UpdatePersonRequest
// @Failure 400,404,500 {object} string
// @Router /persons/{id} [put]
func (h *PersonHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid UUID", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode update request", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), id, req); err != nil {
		h.logger.Error("Failed to update person", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Updated person", zap.String("id", idStr))
	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary Удалить человека по ID
// @Tags persons
// @Param id path string true "UUID человека"
// @Success 204
// @Failure 404,500 {object} string
// @Router /persons/{id} [delete]
func (h *PersonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("Invalid UUID", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Error("Failed to delete person", zap.String("id", idStr), zap.Error(err))
		http.Error(w, "Failed to delete person", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Deleted person", zap.String("id", idStr))
	w.WriteHeader(http.StatusNoContent)
}
