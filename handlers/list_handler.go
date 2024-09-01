package handlers
import(
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"strings"
	"contact-list-api-1/models"
	//"contact-list-api-1/repositories"
	"contact-list-api-1/services"
	"strconv"
	"errors"
	"gorm.io/gorm"
)

type ListHandler struct{
	service services.ListService
}

func NewListHandler(service services.ListService) *ListHandler{
	return &ListHandler{service:service}
}

func (h *ListHandler) GetAllLists(w http.ResponseWriter, r *http.Request){
	queryParams:=r.URL.Query()
	name:=queryParams.Get("name")
	pageStr:=queryParams.Get("page")
	pageSizeStr:=queryParams.Get("pageSize")
	pageNum, err := strconv.Atoi(pageStr)
	
	if err != nil || pageNum <= 0 {
		pageNum=1
	}
	pageSizeNum, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSizeNum <= 0{
		pageSizeNum=10
	}
	lists, err:=h.service.GetAllLists(name,pageNum,pageSizeNum)
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
	parts:=strings.Split(r.URL.Path, "/")
	if len(parts)<4{
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}	
	id:=parts[3]
	uuid, err:=uuid.Parse(id)
	if err!=nil{
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	list, err:=h.service.GetListByUUID(uuid)
	if err != nil {
        	if errors.Is(err, gorm.ErrRecordNotFound) {
            		http.Error(w, "List not found", http.StatusNotFound)
        	} else {
            		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
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
	if list.Name == "" {
		http.Error(w, "Name field is required", http.StatusBadRequest)
		return
	}
	if list.UUID == uuid.Nil {
        	list.UUID = uuid.New()
    	}
	
	if err:=h.service.CreateList(list); err!=nil{
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
func (h *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request){
	parts:=strings.Split(r.URL.Path, "/")
	if len(parts)<4{
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}	
	id:=parts[3]
	uuid, err:=uuid.Parse(id)
	if err!=nil{
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}	
	var list models.List
	err=json.NewDecoder(r.Body).Decode(&list)
	
	if err!=nil{
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	list.UUID=uuid
		
	if err=h.service.UpdateList(list); err!=nil{
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

func (h *ListHandler) DeleteList(w http.ResponseWriter, r *http.Request){
	parts:=strings.Split(r.URL.Path, "/")
	if len(parts)<4{
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}	
	id:=parts[3]
	uuid, err:=uuid.Parse(id)
	if err!=nil{
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	if err:=h.service.DeleteList(uuid); err!=nil{
		if errors.Is(err, gorm.ErrRecordNotFound) {
            		http.Error(w, "List not found", http.StatusNotFound)
        	} else {
            		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
        	}
        	return
    	}
		
	w.WriteHeader(http.StatusNoContent)
}
		
	