
package main
import (
	"contact-list-api-1/config"
	"contact-list-api-1/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"fmt"
)

func main() {
	cfg, err:=config.LoadConfig("config.json")
	if err!=nil{
		log.Fatal("Error loading configuration: ", err)
	}
	dsn:=fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Name)

	db,err:=gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err!=nil{
		log.Fatal("Error connecting to database: " ,err)
	}
	err=db.AutoMigrate(&models.List{}, &models.Contact{})
	if err!=nil{
		log.Fatal(err)
	}

	log.Println("Database connected succesfully")
}
