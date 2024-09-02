package repositories
import (
	"testing"
	
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"github.com/google/uuid"
	
)
func TestContactRepository_GetAll(t *testing.T) {
       	db, cleanup := setTestDB(t)
       	defer cleanup()
            
       	repo := repositories.NewContactRepository(db)

  
       	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
       	}
       	if err := db.Create(&list).Error; err != nil {
		t.Fatalf("Could not create test list: %v", err)
        }

	   
	testContacts := []models.Contact{
		{UUID: uuid.New(), FirstName: "Test", LastName: "Contact", Mobile: "+1234567890", Email: "test.contact@example.com", CountryCode: "USA", ListID: list.ID},
		{UUID: uuid.New(), FirstName: "Test 1", LastName: "Contact 1", Mobile: "+0987654321", Email: "test.test@example.com", CountryCode: "USA", ListID: list.ID},
	}
	for _, contact := range testContacts {
		if err := db.Create(&contact).Error; err != nil {
		    t.Fatalf("Could not create test contact: %v", err)
		}
	}

	t.Run("WithoutFilterAndPagination", func(t *testing.T) {
		contacts, err := repo.GetAll("", "", "", 0, 0)
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(contacts) != len(testContacts) {
		    t.Errorf("Expected %d contacts, got %d", len(testContacts), len(contacts))
		}
	})

	t.Run("WithFilter", func(t *testing.T) {
		contacts, err := repo.GetAll("Contact", "", "", 0, 0)
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(contacts) != 2 {
		    t.Errorf("Expected 2 contacts, got %d", len(contacts))
		}
	})

	t.Run("WithPagination", func(t *testing.T) {
		contacts, err := repo.GetAll("", "", "", 1, 1) 
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(contacts) != 1 {
		    t.Errorf("Expected 1 contact, got %d", len(contacts))
		}
		
	})
	
	
}


func TestContactRepository_Create(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
	}
	if err := db.Create(&list).Error; err != nil {
		t.Fatalf("Could not create test list: %v", err)
	}
	testUUID := uuid.New()
	testContact := models.Contact{
		UUID:        testUUID,
		FirstName:   "Test",
		LastName:    "Contact",
		Mobile:      "+1234567890",
		Email:       "test.contact@example.com",
		CountryCode: "USA",
		ListID:      list.ID, 
	}

	if err := repo.Create(testContact); err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	var createdContact models.Contact
	if err := db.Where("uuid = ?", testUUID).First(&createdContact).Error; err != nil {
		t.Fatalf("Test data not found in database: %v", err)
	}
	t.Run("DuplicateContact", func(t *testing.T) {
	    duplicateContact := testContact
	    err := repo.Create(duplicateContact)
	    if err == nil {
		t.Fatalf("Expected an error due to duplicate email or mobile, got none")
	    }
	})
	
}
func TestContactRepository_GetByUUID(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
	}
	if err := db.Create(&list).Error; err != nil {
		t.Fatalf("Could not create test list: %v", err)
	}
	testUUID := uuid.New()
	testContact := models.Contact{
		UUID:        testUUID,
		FirstName:   "Test",
		LastName:    "Contact",
		Mobile:      "+0987654321",
		Email:       "test.contact@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
	}

	if err := db.Create(&testContact).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	contact, err := repo.GetByUUID(testUUID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if contact == nil || contact.UUID != testUUID {
		t.Errorf("Expected contact with UUID %v, got %v", testUUID, contact)
	
	t.Run("NonExistentUUID", func(t *testing.T) {
	    _, err := repo.GetByUUID(uuid.New())
	    if err == nil {
		t.Fatalf("Expected an error for non-existent UUID, got none")
	    }
	})

}}

func TestContactRepository_Update(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
	}
	if err := db.Create(&list).Error; err != nil {
		t.Fatalf("Could not create test list: %v", err)
	}
	testUUID := uuid.New()
	testContact := models.Contact{
		UUID:        testUUID,
		FirstName:   "Contact",
		LastName:    "Test",
		Mobile:      "+1111111111",
		Email:       "test.test@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
	}

	if err := db.Create(&testContact).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	   
	testContact.LastName = "Contact 1"
	if err := repo.Update(testContact); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var updatedContact models.Contact
	if err := db.Where("uuid = ?", testUUID).First(&updatedContact).Error; err != nil {
		t.Fatalf("Contact not found after update: %v", err)
	}

	if updatedContact.LastName != "Contact 1" {
		t.Errorf("Expected last name 'Contact 1', got '%s'", updatedContact.LastName)
	}
	t.Run("UpdateNonExistentContact", func(t *testing.T) {
		nonExistentUUID := uuid.New() 
		nonExistentContact := models.Contact{
			UUID:        nonExistentUUID,
			FirstName:   "Non",
			LastName:    "Existent",
			Mobile:      "+0000000000",
			Email:       "non.existent@example.com",
			CountryCode: "USA",
			ListID:      list.ID,
		}
		err := repo.Update(nonExistentContact)
		if err == nil {
			t.Fatalf("Expected an error for updating non-existent contact, got none")
		}
	})
	
	
}
func TestContactRepository_Delete(t *testing.T) {
	 db, cleanup := setTestDB(t)
	 defer cleanup()

	 repo := repositories.NewContactRepository(db)
	 list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
	 }
	 if err := db.Create(&list).Error; err != nil {
		t.Fatalf("Could not create test list: %v", err)
	 }
	 testUUID := uuid.New()
	 testContact := models.Contact{
		UUID:        testUUID,
		FirstName:   "Test",
		LastName:    "Contact",
		Mobile:      "+2222222222",
		Email:       "test.test@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
	 }

	 if err := db.Create(&testContact).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	 }

	 if err := repo.Delete(testUUID); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	 }

	 var deletedContact models.Contact
	    result := db.Where("uuid = ?", testUUID).First(&deletedContact)
	    if result.Error == nil {
		t.Errorf("Expected record to be deleted, but it still exists")
	}
	t.Run("DeleteNonExistentUUID", func(t *testing.T) {
	    err := repo.Delete(uuid.New())
	    if err == nil {
		t.Fatalf("Expected an error for deleting non-existent UUID, got none")
	    }
	})
	
}
