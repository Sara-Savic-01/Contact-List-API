package services
import(
	"github.com/google/uuid"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"errors"
)
type ListService interface{
	GetAllLists() ([]models.List, error)
	GetListByUUID(uuid uuid.UUID) (*models.List, error)
	CreateList(list models.List) error
	UpdateList(list models.List) error
	DeleteList(uuid uuid.UUID) error
}

type listService struct{
	repo repositories.ListRepository
}

func NewListService(repo repositories.ListRepository) ListService{
	return &listService{repo:repo}
}
func (s *listService) GetAllLists() ([]models.List, error){
	return s.repo.GetAll()
}
func (s *listService) GetListByUUID(uuid uuid.UUID) (*models.List, error){
		
	return s.repo.GetByUUID(uuid)
}
func (s *listService) CreateList(list models.List) error{
	if list.Name==""{
		return errors.New("Name cannot be empty")
	}       
	
        existingList, err:=s.repo.GetByUUID(list.UUID)
	if err==nil && existingList!=nil{
		return errors.New("List with this UUID already exists")
	}	
	return s.repo.Create(list)
}
func (s *listService) UpdateList(list models.List) error{
	if list.Name==""{
		return errors.New("Name cannot be empty")
	}      
	
        existingList, err:=s.repo.GetByUUID(list.UUID)
	if err==nil && existingList!=nil{
		return errors.New("List with this UUID already exists")
	}	
		
	return s.repo.Update(list)
}
func (s *listService) DeleteList(uuid uuid.UUID) error{
		
	return s.repo.Delete(uuid)
}