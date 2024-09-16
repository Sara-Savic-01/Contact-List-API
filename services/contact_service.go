package services

import (
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"errors"
	"regexp"

	"github.com/google/uuid"
)

type ContactService interface {
	GetAllContacts(name, mobile, email string, page, pageSize int) ([]models.Contact, error)
	GetContactByUUID(uuid uuid.UUID) (*models.Contact, error)
	CreateContact(contact models.Contact) error
	UpdateContact(contact models.Contact) error
	DeleteContact(uuid uuid.UUID) error
}

type contactService struct {
	repo repositories.ContactRepository
}

func NewContactService(repo repositories.ContactRepository) ContactService {
	return &contactService{repo: repo}
}
func (s *contactService) GetAllContacts(name, mobile, email string, page, pageSize int) ([]models.Contact, error) {
	offset := (page - 1) * pageSize
	contacts, err := s.repo.GetAll(name, mobile, email, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}
func (s *contactService) GetContactByUUID(uuid uuid.UUID) (*models.Contact, error) {
	contact, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return contact, nil
}
func (s *contactService) CreateContact(contact models.Contact) error {

	validationErrors := s.validateContact(models.Contact{}, contact, false)
	if validationErrors != nil {
		return validationErrors
	}

	return s.repo.Create(contact)
}
func (s *contactService) UpdateContact(contact models.Contact) error {
	existingContact, err := s.repo.GetByUUID(contact.UUID)
	if err != nil {
		return err
	}
	if existingContact == nil {
		return errors.New("contact not found")
	}

	validationErrors := s.validateContact(*existingContact, contact, true)
	if validationErrors != nil {
		return validationErrors
	}
	return s.repo.Update(contact)
}
func (s *contactService) DeleteContact(uuid uuid.UUID) error {
	existingContact, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return err
	}
	if existingContact == nil {
		return errors.New("contact not found")
	}
	return s.repo.Delete(uuid)
}
func (s *contactService) validateContact(existingContact, contact models.Contact, isUpdate bool) *ValidationErrors {
	var errs []ValidationError

	if isUpdate {

		if contact.Email != "" && contact.Email != existingContact.Email {
			if !isValidEmail(contact.Email) {
				errs = append(errs, ValidationError{Field: "Email", Message: "invalid email format"})
			}
			if isUnique, err := s.isEmailUnique(contact.Email); err != nil {
				errs = append(errs, ValidationError{Field: "Email", Message: "error checking email uniqueness"})
			} else if !isUnique {
				errs = append(errs, ValidationError{Field: "Email", Message: "email already exists"})
			}
		}

		if contact.Mobile != "" && contact.Mobile != existingContact.Mobile {
			if !isValidMobile(contact.Mobile) {
				errs = append(errs, ValidationError{Field: "Mobile", Message: "invalid mobile format"})
			}
			if isUnique, err := s.isMobileUnique(contact.Mobile); err != nil {
				errs = append(errs, ValidationError{Field: "Mobile", Message: "error checking mobile uniqueness"})
			} else if !isUnique {
				errs = append(errs, ValidationError{Field: "Mobile", Message: "mobile already exists"})
			}
		}

		if contact.CountryCode != "" && contact.CountryCode != existingContact.CountryCode {
			if len(contact.CountryCode) != 3 {
				errs = append(errs, ValidationError{Field: "CountryCode", Message: "country code must be exactly 3 characters long"})
			}
		}

	} else {
		if contact.FirstName == "" {
			errs = append(errs, ValidationError{Field: "FirstName", Message: "first name cannot be empty"})

		}
		if contact.LastName == "" {
			errs = append(errs, ValidationError{Field: "LastName", Message: "last name cannot be empty"})

		}
		if contact.Email == "" || !isValidEmail(contact.Email) {
			errs = append(errs, ValidationError{Field: "Email", Message: "invalid email format"})

		}
		if contact.Mobile == "" || !isValidMobile(contact.Mobile) {
			errs = append(errs, ValidationError{Field: "Mobile", Message: "invalid mobile format"})
		}
		if len(contact.CountryCode) != 3 {
			errs = append(errs, ValidationError{Field: "CountryCode", Message: "country code must be exactly 3 characters long"})
		}
		if contact.ListID == 0 {
			errs = append(errs, ValidationError{Field: "ListID", Message: "contact must belong to a list"})
		}
		if isUnique, err := s.isEmailUnique(contact.Email); err != nil {
			errs = append(errs, ValidationError{Field: "Email", Message: "error checking email uniqueness"})
		} else if !isUnique {
			errs = append(errs, ValidationError{Field: "Email", Message: "email already exists"})
		}

		if isUnique, err := s.isMobileUnique(contact.Mobile); err != nil {
			errs = append(errs, ValidationError{Field: "Mobile", Message: "error checking mobile uniqueness"})
		} else if !isUnique {
			errs = append(errs, ValidationError{Field: "Mobile", Message: "mobile already exists"})
		}
		if contact.ListID != 0 {
			exists, err := s.repo.ListExists(contact.ListID)
			if err != nil {
				errs = append(errs, ValidationError{Field: "ListID", Message: "error checking list existence"})
			} else if !exists {
				errs = append(errs, ValidationError{Field: "ListID", Message: "the associated list does not exist"})
			}
		}

	}
	if len(errs) > 0 {
		return NewValidationErrors(errs)
	}
	return nil
}
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
func isValidMobile(mobile string) bool {
	re := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
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
