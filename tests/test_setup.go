package tests
import(
	"testing"
	"contact-list-api-1/config"
	"contact-list-api-1/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"fmt"
)
var db *gorm.DB
func SetupTestDB(t *testing.T) *gorm.DB{
	cfg:=config.LoadTestConfig("/home/osboxes/Documents/contact-list-api-1/tests/config_test.json")
	
	dsn:=fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Name)

	db,err:=gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err!=nil{
		t.Fatalf("Failed to connect to database: %v" ,err)
	}
	err=db.Migrator().DropTable(&models.List{}, &models.Contact{})
	if err!=nil{
		t.Fatalf("Failed to drop tables:%v", err)
	}
	err=db.AutoMigrate(&models.List{}, &models.Contact{})
	if err!=nil{
		t.Fatalf("Failed to migrate tables:%v", err)
	}
	return db
}

func TearDownTestDB(t *testing.T, db *gorm.DB){
	sqlDB, err:=db.DB()
	if err!=nil{
		t.Fatalf("Failed to get db connection: %v",err)
	}
	sqlDB.Close()
}