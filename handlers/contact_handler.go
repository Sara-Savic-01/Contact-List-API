package handlers

import (
	"contact-list-api-1/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	//"contact-list-api-1/repositories"
	"contact-list-api-1/services"
	"errors"

	"gorm.io/gorm"
)

type ContactHandler struct {
	service services.ContactService
}

func NewContactHandler(service services.ContactService) *ContactHandler {
	return &ContactHandler{service: service}
}

func (h *ContactHandler) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	mobile := queryParams.Get("mobile")
	email := queryParams.Get("email")
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
	contacts, err := h.service.GetAllContacts(name, mobile, email, pageNum, pageSizeNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(contacts)

}
func (h *ContactHandler) GetContactByUUID(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	contact, err := h.service.GetContactByUUID(uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Contact not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
}
func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if contact.UUID == uuid.Nil {
		contact.UUID = uuid.New()
	}

	if err := h.service.CreateContact(contact); err != nil {
		var validationErrors *services.ValidationErrors
		if errors.As(err, &validationErrors) {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErrors)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	createdContact, err := h.service.GetContactByUUID(contact.UUID)
	if err != nil {

		http.Error(w, "Failed to retrieve created contact", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdContact)
}

func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	var contact models.Contact
	err = json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	contact.UUID = uuid

	if err = h.service.UpdateContact(contact); err != nil {
		var validationErrors *services.ValidationErrors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Contact not found", http.StatusNotFound)
		} else if errors.As(err, &validationErrors) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErrors)
			return
		} else {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		}
		return

	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteContact(uuid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Contact not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
