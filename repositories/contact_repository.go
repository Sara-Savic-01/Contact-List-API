package repositories
import(
	"contact-list-api-1/models"
	"gorm.io/gorm"
	"errors"
	"github.com/google/uuid"
	"fmt"
)

type ContactRepository interface{
	GetAll(name string, mobile string, email string,limit, offset int) ([]models.Contact, error)
	GetByUUID(uuid uuid.UUID) (*models.Contact, error)
	Create(list models.Contact) error
	Update(list models.Contact) error
	Delete(uuid uuid.UUID) error
	
	
}


type contactRepository struct{
	db *gorm.DB
}


func NewContactRepository(db *gorm.DB) ContactRepository{
	return &contactRepository{db: db}
}
func (c *contactRepository) GetAll(name string, mobile string, email string,limit, offset int) ([]models.Contact, error){
	var contacts []models.Contact
	query:=c.db
	if name != "" {
        	query = query.Where("first_name LIKE ? OR last_name LIKE ?", "%"+name+"%", "%"+name+"%")
    	}
    	if mobile != "" {
        	query = query.Where("mobile LIKE ?", "%"+mobile+"%")
    	}
    	if email != "" {
        	query = query.Where("email LIKE ?", "%"+email+"%")
    	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}
func (c *contactRepository) GetByUUID(uuid uuid.UUID) (*models.Contact, error){
	
		
	var contact models.Contact
	if err:=c.db.Where("uuid =?", uuid).First(&contact).Error; err!=nil{
		return nil, err
	}
	return &contact, nil
}

func (c *contactRepository) Create(contact models.Contact) error{
		
	return c.db.Create(&contact).Error
}
func (c *contactRepository) Update(contact models.Contact) error{
	var existingContact models.Contact
	if err := c.db.Where("uuid = ?", contact.UUID).First(&existingContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Contact with UUID %v does not exist", contact.UUID)
		}
		return err
	}

		
	if err := c.db.Model(&models.Contact{}).Where("uuid = ?", contact.UUID).Updates(contact).Error; err != nil {
		return err
	}
	return nil
}
func (c *contactRepository) Delete(uuid uuid.UUID) error{
	
	var contact models.Contact
	result:=c.db.Where("uuid=?", uuid).First(&contact)
	if result.Error!=nil{
		return result.Error
	}
	result=c.db.Delete(&contact)
	if result.Error!=nil{
		return result.Error
	}
	if result.RowsAffected==0{
		return ErrNotFound
	}	
	return nil
}

