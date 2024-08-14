package repositories
import(
	"contact-list-api-1/models"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type ListRepository interface{
	GetAll() ([]models.List, error)
	GetByUUID(uuid uuid.UUID) (*models.List, error)
	Create(list models.List) error
	Update(list models.List) error
	Delete(uuid uuid.UUID) error
	
}

type listRepository struct{
	db *gorm.DB
}
func NewListRepository(db *gorm.DB) ListRepository{
	return &listRepository{db: db}
}

func (l *listRepository) GetAll() ([]models.List, error){
	var lists []models.List
	if err:=l.db.Find(&lists).Error; err!=nil{
		return nil,err
	}
	return lists, nil
}
func (l *listRepository) GetByUUID(uuid uuid.UUID) (*models.List, error){
	var list models.List
	if err:=l.db.Where("uuid =?", uuid).First(&list).Error; err!=nil{
		return nil, err
	}
	return &list, nil
}

func (l *listRepository) Create(list models.List) error{
	return l.db.Create(&list).Error
}
func (l *listRepository) Update(list models.List) error{
	return l.db.Save(&list).Error
}
func (l *listRepository) Delete(uuid uuid.UUID) error{
	return l.db.Where("uuid=?", uuid).Delete(&models.List{}).Error
}
