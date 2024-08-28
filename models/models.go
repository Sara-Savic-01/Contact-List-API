package models
import(
 "github.com/google/uuid"
 //"gorm.io/gorm"
 
)
type List struct{
	ID uint `gorm:"primaryKey;autoIncrement"`
	UUID uuid.UUID `gorm:"type:char(36); not null;uniqueIndex" json:"uuid":` 
	Name string `gorm:"type:varchar(255);not null" json:"name"`
	
}

type Contact struct{
	ID uint `gorm:"primaryKey;autoIncrement"`
	UUID uuid.UUID `gorm:"type:char(36); not null;uniqueIndex" json:"uuid:` 
	FirstName string `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName string `gorm:"type:varchar(255);not null" json:"last_name"`
	Mobile string `gorm:"type:varchar(20);not null" json:"mobile"`
	Email string `gorm:"type:varchar(255); not null ; uniqueIndex" json:"email"`
	CountryCode string `gorm:"type:varchar(3);not null" json:"country_code"`
	ListID uint `gorm:"not null" json:"list_id"`
	
}

