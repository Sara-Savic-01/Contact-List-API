package main

import (
	"contact-list-api-1/config"
	"contact-list-api-1/handlers"
	middleware "contact-list-api-1/middlewares"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"contact-list-api-1/services"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	err = db.AutoMigrate(&models.List{}, &models.Contact{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connected succesfully")

	listRepo := repositories.NewListRepository(db)

	contactRepo := repositories.NewContactRepository(db)
	listService := services.NewListService(listRepo)
	contactService := services.NewContactService(contactRepo)

	listHandler := handlers.NewListHandler(listService)
	contactHandler := handlers.NewContactHandler(contactService)

	http.Handle("GET /lists", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(listHandler.GetAllLists)))
	http.Handle("GET /lists/{uuid}", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(listHandler.GetListByUUID)))
	http.Handle("POST /lists", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(listHandler.CreateList)))
	http.Handle("PUT /lists/{uuid}", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(listHandler.UpdateList)))
	http.Handle("DELETE /lists/{uuid}", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(listHandler.DeleteList)))

	http.Handle("GET /contacts", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(contactHandler.GetAllContacts)))
	http.Handle("GET /contacts/{uuid}", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(contactHandler.GetContactByUUID)))
	http.Handle("POST /contacts", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(contactHandler.CreateContact)))
	http.Handle("PUT /contacts/{uuid}", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(contactHandler.UpdateContact)))
	http.Handle("DELETE /contacts/{uuid}", middleware.AuthMiddleware(cfg.AuthToken, http.HandlerFunc(contactHandler.DeleteContact)))

	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
