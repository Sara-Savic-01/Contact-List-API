package repositories
import(
	"contact-list-api-1/models"
	"gorm.io/gorm"
	"errors"
	"github.com/google/uuid"
)

type ListRepository interface{
	
	GetAll(name string, limit, offset int) ([]models.List, error)
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

func (l *listRepository) GetAll(name string, limit, offset int) ([]models.List, error){
	if limit < 0 || offset < 0 {
        	return nil, errors.New("limit and offset must be non-negative")
    	}
	var lists []models.List
	query:=l.db
	
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&lists).Error; err != nil {
		return nil, err
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
	
	if err := l.db.Model(&models.List{}).Where("uuid = ?", list.UUID).Updates(list).Error; err != nil {
		return err
	}
	return nil
}
func (l *listRepository) Delete(uuid uuid.UUID) error{
	
	var list models.List
	result:=l.db.Where("uuid = ?", uuid).First(&list)
	if result.Error!=nil{
		return result.Error
	}
	l.db.Where("list_id = ?", list.ID).Delete(&models.Contact{})
	result=l.db.Delete(&list)
	if result.Error!=nil{
		return result.Error
	}
	if result.RowsAffected==0{
		return ErrNotFound
	}	
	return nil
}
