package services
import(
	"github.com/google/uuid"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
)
type ContactService interface{
	GetAllContacts() ([]models.Contact, error)
	GetContactByUUID(uuid uuid.UUID) (*models.Contact, error)
	CreateContact(contact models.Contact) error
	UpdateContact(contact models.Contact) error
	DeleteContact(uuid uuid.UUID) error
	
}


type contactService struct{
	repo repositories.ContactRepository
}

func NewContactService(repo repositories.ContactRepository) ContactService{
	return &contactService{repo:repo}
}
func (s *contactService) GetAllContacts() ([]models.Contacts, error){
	return s.repo.GetAll()
}
func (s *contactService) GetContactByUUID(uuid uuid.UUID) (*models.Contact, error){
	return s.repo.GetByUUID(uuid)
}
func (s *contactService) CreateContact(contact models.Contact) error{
	return s.rep.Create(contact)
}
func (s *contactService) UpdateContact(contact models.Contact) error{
	return s.repo.Update(contact)
}
func (s *contactService) DeleteContact(uuid uuid.UUID) error{
	return s.repo.Delete(uuid)
}
