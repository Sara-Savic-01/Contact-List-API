package handlers
import(
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"regexp"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"contact-list-api-1/services"
)

type ContactHandler struct{
	service services.ContactService
}
func NewContactHandler(service services.ContactService) *ContactHandler{
	return &ContactHandler{service:service}
}
func (h *ContactHandler) GetAllContacts(w http.ResponseWriter, r *http.Request){
	contacts, err:=h.service.GetAllContacts()
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data,err:=json.Marshal(contacts)
	if err!=nil{
		http.Error(w, "Failed to load data ", http.StatusInternalServerError)
		return
	}
	w.Write(data)

}
func (h *ContactHandler) GetContactByUUID(w http.ResponseWriter, r *http.Request){
	id:=r.URL.Query().Get("uuid")
	uuid, err:=uuid.Parse(id)
	if err!=nil{
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	contact, err:=h.service.GetContactByUUID(uuid)
	if err!=nil{
		if err==repositories.ErrNotFound{
			http.Error(w, "Contact not found", http.StatusNotFound)
		}else{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
		w.Header().Set("Content-Type", "application/json")
	data, err:=json.Marshal(contact)
	if err!=nil{
		http.Error(w, "Failed to load data ", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request){
	var contact models.Contact
	err:=json.NewDecoder(r.Body).Decode(&contact)
	if err!=nil{
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	if contact.Email!="" && !isValidEmail(contact.Email){
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}
	if contact.Mobile!="" && !isValidMobile(contact.Mobile){
		http.Error(w, "invalid mobile format", http.StatusBadRequest)
		return
	}
	
	
	if err:=h.service.CreateContact(contact);err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request){
	var contact models.Contact
	err:=json.NewDecoder(r.Body).Decode(&contact)
	if err!=nil{
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	
	if contact.Email!="" && !isValidEmail(contact.Email){
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}
	if contact.Mobile!="" && !isValidMobile(contact.Mobile){
		http.Error(w, "invalid mobile format", http.StatusBadRequest)
		return
	}
	if err:=h.service.UpdateContact(contact);err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request){
	id:=r.URL.Query().Get("uuid")
	uuid, err:=uuid.Parse(id)
	if err!=nil{
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	if err:=h.service.DeleteContact(uuid); err!=nil{
		if err==repositories.ErrNotFound{
			http.Error(w, "List not found", http.StatusNotFound)
		}else{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func isValidEmail(email string) bool{
	re:=regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
func isValidMobile(mobile string) bool{
	re:=regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return re.MatchString(mobile)
}
