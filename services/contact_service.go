package services
import(
	"github.com/google/uuid"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"errors"
	"regexp"
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
func (s *contactService) GetAllContacts() ([]models.Contact, error){
	return s.repo.GetAll()
}
func (s *contactService) GetContactByUUID(uuid uuid.UUID) (*models.Contact, error){
		
	return s.repo.GetByUUID(uuid)
}
func (s *contactService) CreateContact(contact models.Contact) error{
	
	if contact.Email!=""&&!isValidEmail(contact.Email){
		return errors.New("Invalid email format")
	}	
	if contact.Mobile!=""&&!isValidMobile(contact.Mobile){
		return errors.New("Invalid mobile format")
	}
	
	if contact.ListID==0{
		return errors.New("Contact must belong to a list")
	}
	return s.repo.Create(contact)
}
func (s *contactService) UpdateContact(contact models.Contact) error{
	
	if contact.Email!=""&&!isValidEmail(contact.Email){
		return errors.New("Invalid email format")
	}	
	if contact.Mobile!=""&&!isValidMobile(contact.Mobile){
		return errors.New("Invalid mobile format")
	}
	
	if contact.ListID==0{
		return errors.New("Contact must belong to a list")
	}
	existingContact, err:=s.repo.GetByUUID(contact.UUID)
	if err!=nil || existingContact==nil{
		return errors.New("Contact not found")
	}
	return s.repo.Update(contact)
}
func (s *contactService) DeleteContact(uuid uuid.UUID) error{
	
	
	
	return s.repo.Delete(uuid)
}

func isValidEmail(email string) bool{
	re:=regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
func isValidMobile(mobile string) bool{
	re:=regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return re.MatchString(mobile)
}
