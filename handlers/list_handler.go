package handlers

import (
	"contact-list-api-1/models"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	//"contact-list-api-1/repositories"
	"contact-list-api-1/services"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type ListHandler struct {
	service services.ListService
}

func NewListHandler(service services.ListService) *ListHandler {
	return &ListHandler{service: service}
}

func (h *ListHandler) GetAllLists(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	pageStr := queryParams.Get("page")
	pageSizeStr := queryParams.Get("pageSize")
	pageNum, err := strconv.Atoi(pageStr)

	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSizeNum, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}
	lists, err := h.service.GetAllLists(name, pageNum, pageSizeNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lists)

}

func (h *ListHandler) GetListByUUID(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("uuid")

	uuid, err := uuid.Parse(id)

	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	list, err := h.service.GetListByUUID(uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "List not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)

}

func (h *ListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	var list models.List
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if list.UUID == uuid.Nil {
		list.UUID = uuid.New()
	}

	if err := h.service.CreateList(list); err != nil {
		var validationErrors *services.ValidationErrors
		if errors.As(err, &validationErrors) {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErrors)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	createdList, err := h.service.GetListByUUID(list.UUID)
	if err != nil {
		http.Error(w, "Failed to retrieve created list", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdList)
}
func (h *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	var list models.List
	err = json.NewDecoder(r.Body).Decode(&list)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	list.UUID = uuid

	if err = h.service.UpdateList(list); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "List not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(list)
}

func (h *ListHandler) DeleteList(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteList(uuid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "List not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
