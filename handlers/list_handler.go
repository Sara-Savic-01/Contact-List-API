package handlers
import(
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"contact-list-api-1/services"
)

type ListHandler struct{
	service services.ListService
}

func NewListHandler(service services.ListService) *ListHandler{
	return &ListHandler{service:service}
}

func (h *ListHandler) GetAllLists(w http.ResponseWriter, r *http.Request){
	lists, err:=h.service.GetAllLists()
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, err:=json.Marshal(lists)
	if err!=nil{
		http.Error(w, "Failed to load data ", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *ListHandler) GetListByUUID(w http.ResponseWriter, r *http.Request){
	id:=r.URL.Query().Get("uuid")
	uuid, err:=uuid.Parse(id)
	if err!=nil{
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	list, err:=h.service.GetListByUUID(uuid)
	if err!=nil{
		if err==repositories.ErrNotFound{
			http.Error(w, "List not found", http.StatusNotFound)
		}else{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
		w.Header().Set("Content-Type", "application/json")
	data, err:=json.Marshal(list)
	if err!=nil{
		http.Error(w, "Failed to load data ", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *ListHandler) CreateList(w http.ResponseWriter, r *http.Request){
	var list models.List
	err:=json.NewDecoder(r.Body).Decode(&list)
	if err!=nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if list.Name==""{
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}
	
	if err:=h.service.CreateList(list); err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(list)
}
func (h *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request){
	var list models.List
	err:=json.NewDecoder(r.Body).Decode(&list)
	
	if err!=nil{
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err:=h.service.UpdateList(list); err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(list)
}

func (h *ListHandler) DeleteList(w http.ResponseWriter, r *http.Request){
	id:=r.URL.Query().Get("uuid")
	uuid, err:=uuid.Parse(id)
	if err!=nil{
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	if err:=h.service.DeleteList(uuid); err!=nil{
		if err==repositories.ErrNotFound{
			http.Error(w, "List not found", http.StatusNotFound)
		}else{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
		
	