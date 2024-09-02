package services
import (
	"testing"
	"contact-list-api-1/services"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"github.com/google/uuid"
	
)

func TestContactService_GetAllContacts(t *testing.T) {
    	db, cleanup := setTestDB(t)
    	defer cleanup()

    	repo := repositories.NewContactRepository(db)
    	service := services.NewContactService(repo)

    	list := models.List{
        	UUID: uuid.New(),
        	Name: "Test List",
    	}
    	if err := db.Create(&list).Error; err != nil {
        	t.Fatalf("Could not create test list: %v", err)
    	}

    	contacts := []models.Contact{
		{
		    UUID:        uuid.New(),
		    FirstName:   "Sara",
		    LastName:    "Test",
		    Mobile:      "+1234567890",
		    Email:       "test.test@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		},
		{
		    UUID:        uuid.New(),
		    FirstName:   "Test 1",
		    LastName:    "Test",
		    Mobile:      "+0987654321",
		    Email:       "contact.contact@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		},
		{
		    UUID:        uuid.New(),
		    FirstName:   "Test",
		    LastName:    "Test",
		    Mobile:      "+1122334455",
		    Email:       "test1.test1@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		},
    	}
    	for _, contact := range contacts {
        	if err := db.Create(&contact).Error; err != nil {
            		t.Fatalf("Could not create test contact: %v", err)
        	}
    	}

    	t.Run("WithoutFilterAndPagination", func(t *testing.T) {
		results, err := service.GetAllContacts("", "", "", 1, 10) 
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(results) != len(contacts) {
		    t.Errorf("Expected %d contacts, got %d", len(contacts), len(results))
		}
    	})

    	t.Run("WithFilter", func(t *testing.T) {
	
		results, err := service.GetAllContacts("Sara", "", "", 1, 10)
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(results) != 1 {
		    t.Errorf("Expected 1 contact with name 'Sara', got %d", len(results))
		}
		if results[0].FirstName != "Sara" {
		    t.Errorf("Expected contact with name 'Sara', got %s", results[0].FirstName)
		}
    	})

    	t.Run("WithPagination", func(t *testing.T) {
		results, err := service.GetAllContacts("", "", "", 1, 1) 
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(results) != 1 {
		    t.Errorf("Expected 1 contact on page 1 with page size 1, got %d", len(results))
		}
        
    	})
    

}

func TestContactService_GetContactByUUID(t *testing.T) {
    	db, cleanup := setTestDB(t)
    	defer cleanup()

    	repo := repositories.NewContactRepository(db)
    	service := services.NewContactService(repo)

    	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
    	}
    	db.Create(&list)

    	testUUID := uuid.New()
    	contact := models.Contact{
		UUID:        testUUID,
		FirstName:   "Test",
		LastName:    "Test",
		Mobile:      "+1122334455",
		Email:       "test1.test1@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
    	}
    	db.Create(&contact)

    	result, err := service.GetContactByUUID(testUUID)
    	if err != nil {
        	t.Fatalf("Expected no error, got %v", err)
    	}
    	if result.UUID != testUUID {
        	t.Errorf("Expected UUID %v, got %v", testUUID, result.UUID)
    	}
    	t.Run("NonExistentUUID", func(t *testing.T) {
		nonExistentUUID := uuid.New()
		result, err := service.GetContactByUUID(nonExistentUUID)
		if err == nil {
		    t.Fatalf("Expected an error, got %v", err)
		}
		if result != nil {
		    t.Errorf("Expected no contact, got %v", result)
		}
    	})
}

func TestContactService_CreateContact(t *testing.T) {
    	db, cleanup := setTestDB(t)
    	defer cleanup()

    	repo := repositories.NewContactRepository(db)
    	service := services.NewContactService(repo)

    	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
    	}
    	db.Create(&list)

    	newContact := models.Contact{
		UUID:        uuid.New(),
		FirstName:   "Test",
		LastName:    "Test",
		Mobile:      "+1122334455",
		Email:       "test1.test1@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
    	}

    	err := service.CreateContact(newContact)
    	if err != nil {
        	t.Fatalf("Expected no error, got %v", err)
    	}

    	var contact models.Contact
    	db.Where("uuid = ?", newContact.UUID).First(&contact)
    	if contact.Email != newContact.Email {
        	t.Errorf("Expected email %v, got %v", newContact.Email, contact.Email)
    	}
	t.Run("CreateContact_InvalidListID", func(t *testing.T) {
		contact := models.Contact{
		    FirstName:   "Test",
		    LastName:    "Test",
		    Mobile:      "+987654321",
		    Email:       "test.test@example.com",
		    CountryCode: "USA",
		    ListID:      111111,
		}

		err := service.CreateContact(contact)
		if err == nil {
		    t.Fatal("Expected error, got none")
		}

		
    	})
    	t.Run("InvalidEmailFormat", func(t *testing.T) {
		list := models.List{
		    UUID: uuid.New(),
		    Name: "Test List",
		}
		db.Create(&list)
		invalidEmailContact := models.Contact{
		    UUID:        uuid.New(),
		    FirstName:   "Test",
		    LastName:    "Test",
		    Mobile:      "+1122334455",
		    Email:       "invalid-email",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		err := service.CreateContact(invalidEmailContact)
		if err == nil {
		    t.Fatal("Expected error for invalid email format, got none")
		}
    	})

    	t.Run("DuplicateEmail", func(t *testing.T) {
		list := models.List{
		    UUID: uuid.New(),
		    Name: "Test List",
		}
		db.Create(&list)
		existingContact := models.Contact{
		    UUID:        uuid.New(),
		    FirstName:   "Existing",
		    LastName:    "Contact",
		    Mobile:      "+1234567890",
		    Email:       "duplicate@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		if err := db.Create(&existingContact).Error; err != nil {
		    t.Fatalf("Could not create existing contact: %v", err)
		}

		newContact := models.Contact{
		    UUID:        uuid.New(),
		    FirstName:   "New",
		    LastName:    "Contact",
		    Mobile:      "+0987654321",
		    Email:       "duplicate@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		err := service.CreateContact(newContact)
		if err == nil {
		    t.Fatal("Expected error for duplicate email, got none")
		}
    	})
    	t.Run("InvalidMobileFormat", func(t *testing.T) {
		list := models.List{
		    UUID: uuid.New(),
		    Name: "Test List",
		}
		db.Create(&list)
		invalidMobileContact := models.Contact{
		    UUID:        uuid.New(),
		    FirstName:   "Test",
		    LastName:    "Test",
		    Mobile:      "1122334455",
		    Email:       "invalid@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		err := service.CreateContact(invalidMobileContact)
		if err == nil {
		    t.Fatal("Expected error for invalid email format, got none")
		}
    	})

    	t.Run("DuplicateMobile", func(t *testing.T) {
		list := models.List{
		    UUID: uuid.New(),
		    Name: "Test List",
		}
		db.Create(&list)
		existingContact := models.Contact{
		    UUID:        uuid.New(),
		    FirstName:   "Existing",
		    LastName:    "Contact",
		    Mobile:      "+1234567890",
		    Email:       "existing@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		if err := db.Create(&existingContact).Error; err != nil {
		    t.Fatalf("Could not create existing contact: %v", err)
		}

		newContact := models.Contact{
		    UUID:        uuid.New(),
		    FirstName:   "New",
		    LastName:    "Contact",
		    Mobile:      "+1234567890",
		    Email:       "test@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		err := service.CreateContact(newContact)
		if err == nil {
		    t.Fatal("Expected error for duplicate mobile, got none")
		}
	    })
    	t.Run("EmptyFieldsAndInvalidCountryCode", func(t *testing.T) {
		
			invalidContact := models.Contact{
			    UUID:        uuid.New(),
			    FirstName:   "",
			    LastName:    "",
			    Mobile:      "",
			    Email:       "",
			    CountryCode: "US",
				
			}
			err := service.CreateContact(invalidContact)
			if err == nil {
			    t.Fatal("Expected error for invalid email format, got none")
			}
	      })
		
}

func TestContactService_UpdateContact(t *testing.T) {
    	db, cleanup := setTestDB(t)
    	defer cleanup()

    	repo := repositories.NewContactRepository(db)
    	service := services.NewContactService(repo)

    	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
    	}
    	db.Create(&list)

    	existingUUID := uuid.New()
    	contact := models.Contact{
		UUID:        existingUUID,
		FirstName:   "Test",
		LastName:    "Test",
		Mobile:      "+1122334455",
		Email:       "test1.test1@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
    	}
    	db.Create(&contact)

    	contact.Email = "test.new@example.com"
    	err := service.UpdateContact(contact)
    	if err != nil {
        	t.Fatalf("Expected no error, got %v", err)
    	}

    	var updatedContact models.Contact
    	db.Where("uuid = ?", existingUUID).First(&updatedContact)
    	if updatedContact.Email != contact.Email {
        	t.Errorf("Expected email %v, got %v", contact.Email, updatedContact.Email)
    	}
	t.Run("UpdateContact_InvalidListID", func(t *testing.T) {
        	
		invalidContact := models.Contact{
		    UUID:        contact.UUID,
		    FirstName:   "Test",
		    LastName:    "Test",
		    Mobile:      "+1122334455",
		    Email:       "invalid-email",
		    CountryCode: "USA",
		    ListID:      111111,
		}

        err := service.UpdateContact(invalidContact)
        if err == nil {
            t.Fatal("Expected error, got none")
        }

        
	})
    	t.Run("InvalidEmailFormat", func(t *testing.T) {
		
		invalidEmailContact := models.Contact{
		    UUID:        contact.UUID,
		    FirstName:   "Test",
		    LastName:    "Test",
		    Mobile:      "+1122334455",
		    Email:       "invalid-email",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		err := service.UpdateContact(invalidEmailContact)
		if err == nil {
		    t.Fatal("Expected error for invalid email format, got none")
		}
    	})

    	t.Run("InvalidMobileFormat", func(t *testing.T) {
		
		invalidMobileContact := models.Contact{
		    UUID:        contact.UUID,
		    FirstName:   "Test",
		    LastName:    "Test",
		    Mobile:      "1122334455",
		    Email:       "invalid@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		err := service.UpdateContact(invalidMobileContact)
		if err == nil {
		    t.Fatal("Expected error for invalid email format, got none")
		}
    	})
    	t.Run("NonExistentUUID", func(t *testing.T) {
		list := models.List{
		    UUID: uuid.New(),
		    Name: "Test List",
		}
		nonExistentUUID := uuid.New()
		contact := models.Contact{
		    UUID:        nonExistentUUID,
		    FirstName:   "Test",
		    LastName:    "Test",
		    Mobile:      "+1122334455",
		    Email:       "test1.test1@example.com",
		    CountryCode: "USA",
		    ListID:      list.ID,
		}
		err := service.UpdateContact(contact)
		if err == nil {
		    t.Fatal("Expected error for updating non-existent contact, got none")
		}
    	})
    	t.Run("InvalidCountryCode", func(t *testing.T) {
		
		invalidContact := models.Contact{
		    UUID:        contact.UUID,
		    FirstName:   "",
		    LastName:    "",
		    Mobile:      "",
		    Email:       "",
		    CountryCode: "US",
		        
		}
		err := service.UpdateContact(invalidContact)
		if err == nil {
		    t.Fatal("Expected error for invalid email format, got none")
		}
      	})

}

func TestContactService_DeleteContact(t *testing.T) {
    	db, cleanup := setTestDB(t)
    	defer cleanup()

    	repo := repositories.NewContactRepository(db)
    	service := services.NewContactService(repo)

    	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
    	}
    	db.Create(&list)

    	testUUID := uuid.New()
    	contact := models.Contact{
		UUID:        testUUID,
		FirstName:   "Test",
		LastName:    "Test",
		Mobile:      "+1122334455",
		Email:       "test1.test1@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
    	}
    	db.Create(&contact)

    	err := service.DeleteContact(testUUID)
    	if err != nil {
        	t.Fatalf("Expected no error, got %v", err)
    	}

    	var deletedContact models.Contact
    	if result := db.Where("uuid = ?", testUUID).First(&deletedContact); result.Error == nil {
        	t.Errorf("Expected record to be deleted, but it still exists")
    	}
    
     	t.Run("DeleteAlreadyDeleted", func(t *testing.T) {
		    validUUID := uuid.New()
		    contact := models.Contact{
			UUID:        testUUID,
			FirstName:   "Test",
			LastName:    "Test",
			Mobile:      "+1122334455",
			Email:       "test1.test1@example.com",
			CountryCode: "USA",
			ListID:      list.ID,
		    }
		    db.Create(&contact)
		    db.Delete(&contact) 

		    err := service.DeleteContact(validUUID)
		    if err == nil {
			t.Fatalf("Expected error when deleting already deleted list, got %v", err)
		    }

		    var deletedContact models.Contact
		    if result := db.Where("uuid = ?", validUUID).First(&deletedContact); result.Error == nil {
			t.Errorf("Expected record to be deleted, but it still exists")
		    }
    	})
    	t.Run("NonExistentContact", func(t *testing.T) {
		nonExistentUUID := uuid.New()
		err := service.DeleteContact(nonExistentUUID)
		if err == nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		
		var contact models.Contact
		if result := db.Where("uuid = ?", nonExistentUUID).First(&contact); result.Error == nil {
		    t.Errorf("Expected contact to not exist, but it was found")
		}
	    })
}
