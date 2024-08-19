package repositories
import(
	"contact-list-api-1/models"
	"gorm.io/gorm"
	
	"github.com/google/uuid"
)

type ContactRepository interface{
	GetAll() ([]models.Contact, error)
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

func (c *contactRepository) GetAll() ([]models.Contact, error){
	var contacts []models.Contact
	if err:=c.db.Find(&contacts).Error; err!=nil{
		return nil,err
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

