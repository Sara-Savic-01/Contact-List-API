package services
import(
	"github.com/google/uuid"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"errors"
	"regexp"
	
)
type ContactService interface{
	GetAllContacts(name, mobile, email string, page,pageSize int) ([]models.Contact, error)
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
func (s *contactService) GetAllContacts(name, mobile, email string, page,pageSize int) ([]models.Contact, error){	
	offset:=(page - 1)*pageSize
	contacts, err:=	s.repo.GetAll(name,mobile,email,pageSize,offset)
	if err != nil {
		return nil, err
	}
	return contacts,nil
}
func (s *contactService) GetContactByUUID(uuid uuid.UUID) (*models.Contact, error){
	contact, err:=s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return contact, nil
}
func (s *contactService) CreateContact(contact models.Contact) error{
	exists, err := s.repo.ListExists(contact.ListID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("The associated list does not exist")
	}

	if contact.Email!=""&&!isValidEmail(contact.Email){
		return errors.New("Invalid email format")
	}	
	if contact.Mobile!=""&&!isValidMobile(contact.Mobile){
		return errors.New("Invalid mobile format")
	}
	if contact.FirstName=="" || contact.LastName==""{
		return errors.New("First and last name cannot be empty")
	}
	if isUnique, err := s.isEmailUnique(contact.Email); err != nil {
        	return err
    	} else if !isUnique {
        	return errors.New("Email already exists")
    	}

    	if isUnique, err := s.isMobileUnique(contact.Mobile); err != nil {
        	return err
    	} else if !isUnique {
        	return errors.New("Mobile already exists")
    	}
	if len(contact.CountryCode)!=3{
		return errors.New("Country code must be exactly 3 characters long")
	}
	if contact.UUID == uuid.Nil {
        	contact.UUID = uuid.New()
    	}
	if contact.ListID==0{
		return errors.New("Contact must belong to a list")
	}
	
	return s.repo.Create(contact)
}
func (s *contactService) UpdateContact(contact models.Contact) error{
	/*exists, err := s.repo.ListExists(contact.ListID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("The associated list does not exist")
	}*/
	existingContact, err:=s.repo.GetByUUID(contact.UUID)
	if err!=nil{
		return err
	}
	if existingContact==nil{
		return errors.New("Contact not found")
	}
	if contact.Email != "" && contact.Email != existingContact.Email {
		if !isValidEmail(contact.Email){
			return errors.New("Invalid email format")
		}	
		existingContact.Email = contact.Email
	}
	if contact.Mobile != "" && contact.Mobile != existingContact.Mobile {
		if !isValidMobile(contact.Mobile){
			return errors.New("Invalid mobile format")
		}
		existingContact.Mobile = contact.Mobile
	}
	if contact.CountryCode != "" && contact.CountryCode != existingContact.CountryCode {
		if len(contact.CountryCode)!=3{
			return errors.New("Country code must be exactly 3 characters long")
		}
		existingContact.CountryCode = contact.CountryCode
	}
	return s.repo.Update(contact)
}
func (s *contactService) DeleteContact(uuid uuid.UUID) error{
	existingContact, err:=s.repo.GetByUUID(uuid)
	if err!=nil{
		return err
	}
	if existingContact==nil{
		return errors.New("Contact not found")
	}
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
func (s *contactService) isEmailUnique(email string) (bool, error) {
    	contacts, err := s.repo.GetAll("", "", email, 1, 0)
    	if err != nil {
        	return false, err
    	}
    	return len(contacts) == 0, nil
}

func (s *contactService) isMobileUnique(mobile string) (bool, error) {
    	contacts, err := s.repo.GetAll("", mobile, "", 1, 0)
    	if err != nil {
        	return false, err
    	}
    	return len(contacts) == 0, nil
}
