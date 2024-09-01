package services
import(
	"github.com/google/uuid"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"errors"
)
type ListService interface{
	GetAllLists(name string, page, pageSize int) ([]models.List, error)
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
func (s *listService) GetAllLists(name string, page, pageSize int) ([]models.List, error){
		
	offset:=(page -1)*pageSize	
	lists, err := s.repo.GetAll(name,pageSize,offset)
	if err != nil {
		return nil, err
	}
	return lists, nil
}
func (s *listService) GetListByUUID(uuid uuid.UUID) (*models.List, error){
	
	list, err := s.repo.GetByUUID(uuid)
    	if err != nil {
        	return nil, err
    	}
   	return list, nil
}
func (s *listService) CreateList(list models.List) error{
	if list.Name==""{
		return errors.New("Name cannot be empty")
	}
	if list.UUID == uuid.Nil {
        	list.UUID = uuid.New()
    	}       
		
	
	return s.repo.Create(list)
}
func (s *listService) UpdateList(list models.List) error{
	existingList, err:=s.repo.GetByUUID(list.UUID)
	if err!=nil{
		return err
	}
	if existingList==nil{
		return errors.New("List not found")
	}
	return s.repo.Update(list)
}
func (s *listService) DeleteList(uuid uuid.UUID) error{
	existingList, err:=s.repo.GetByUUID(uuid)
	if err!=nil{
		return err
	}
	if existingList==nil{
		return errors.New("List not found")
	}

	return s.repo.Delete(uuid)
}