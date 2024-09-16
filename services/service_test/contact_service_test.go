package services

import (
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"contact-list-api-1/services"
	"testing"

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
			Mobile:      "+1987654321",
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
	testCases := []struct {
		name          string
		filterName    string
		filterMobile  string
		filterEmail   string
		page          int
		pageSize      int
		expectedCount int
		expectedFirst string
	}{
		{
			name:          "WithoutFilterAndWithDefaultPagination",
			filterName:    "",
			filterMobile:  "",
			filterEmail:   "",
			page:          1,
			pageSize:      10,
			expectedCount: len(contacts),
		},
		{
			name:          "WithFilter",
			filterName:    "Sara",
			filterMobile:  "",
			filterEmail:   "",
			page:          1,
			pageSize:      10,
			expectedCount: 1,
			expectedFirst: "Sara",
		},
		{
			name:          "WithPagination",
			filterName:    "",
			filterMobile:  "",
			filterEmail:   "",
			page:          1,
			pageSize:      1,
			expectedCount: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.GetAllContacts(tt.filterName, tt.filterMobile, tt.filterEmail, tt.page, tt.pageSize)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d contacts, got %d", tt.expectedCount, len(results))
			}
			if tt.expectedFirst != "" && len(results) > 0 && results[0].FirstName != tt.expectedFirst {
				t.Errorf("Expected contact with first name '%s', got %s", tt.expectedFirst, results[0].FirstName)
			}
		})
	}

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

	testCases := []struct {
		name          string
		contactUUID   uuid.UUID
		expectedError bool
		expectedUUID  uuid.UUID
	}{
		{
			name:          "ExistingUUID",
			contactUUID:   testUUID,
			expectedError: false,
			expectedUUID:  testUUID,
		},
		{
			name:          "NonExistentUUID",
			contactUUID:   uuid.New(),
			expectedError: true,
			expectedUUID:  uuid.Nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetContactByUUID(tt.contactUUID)
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got %v", tt.expectedError, err)
			}
			if tt.expectedError && result != nil {
				t.Errorf("Expected no contact, got %v", result)
			} else if !tt.expectedError && result.UUID != tt.expectedUUID {
				t.Errorf("Expected UUID %v, got %v", tt.expectedUUID, result.UUID)
			}
		})
	}
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
	if err := db.Create(&list).Error; err != nil {
		t.Fatalf("Failed to create test list: %v", err)
	}

	testCases := []struct {
		name             string
		contact          models.Contact
		existingContacts []models.Contact
		expectedError    bool
	}{
		{
			name: "ValidContact",
			contact: models.Contact{
				UUID:        uuid.New(),
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "+1122334455",
				Email:       "test1.test1@example.com",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			expectedError: false,
		},
		{
			name: "InvalidListID",
			contact: models.Contact{
				UUID:        uuid.New(),
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "+987654321",
				Email:       "test.test@example.com",
				CountryCode: "USA",
				ListID:      111111,
			},
			expectedError: true,
		},
		{
			name: "InvalidEmailFormat",
			contact: models.Contact{
				UUID:        uuid.New(),
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "+1122334455",
				Email:       "invalid-email",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			expectedError: true,
		},
		{
			name: "DuplicateEmail",
			contact: models.Contact{
				UUID:        uuid.New(),
				FirstName:   "New",
				LastName:    "Contact",
				Mobile:      "+1987654321",
				Email:       "duplicate@example.com",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			existingContacts: []models.Contact{
				{
					UUID:        uuid.New(),
					FirstName:   "Existing",
					LastName:    "Contact",
					Mobile:      "+1234567890",
					Email:       "duplicate@example.com",
					CountryCode: "USA",
					ListID:      list.ID,
				},
			},
			expectedError: true,
		},
		{
			name: "InvalidMobileFormat",
			contact: models.Contact{
				UUID:        uuid.New(),
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "1122334455",
				Email:       "valid@example.com",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			expectedError: true,
		},
		{
			name: "DuplicateMobile",
			contact: models.Contact{
				UUID:        uuid.New(),
				FirstName:   "New",
				LastName:    "Contact",
				Mobile:      "+1234567890",
				Email:       "unique@example.com",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			existingContacts: []models.Contact{
				{
					UUID:        uuid.New(),
					FirstName:   "Existing",
					LastName:    "Contact",
					Mobile:      "+1234567890",
					Email:       "existing@example.com",
					CountryCode: "USA",
					ListID:      list.ID,
				},
			},
			expectedError: true,
		},
		{
			name: "EmptyFieldsAndInvalidCountryCode",
			contact: models.Contact{
				UUID:        uuid.New(),
				FirstName:   "",
				LastName:    "",
				Mobile:      "",
				Email:       "",
				CountryCode: "US",
			},
			expectedError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			for _, existingContact := range tt.existingContacts {
				if err := db.Create(&existingContact).Error; err != nil {
					t.Fatalf("Failed to create existing contact for test: %v", err)
				}
			}

			err := service.CreateContact(tt.contact)
			if tt.expectedError {
				if err == nil {
					t.Fatal("Expected error, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				var contact models.Contact
				if err := db.Where("uuid = ?", tt.contact.UUID).First(&contact).Error; err != nil {
					t.Fatalf("Failed to find contact: %v", err)
				}
				if contact.UUID == uuid.Nil {
					t.Errorf("Expected UUID to be generated, got nil")
				}
				if contact.Email != tt.contact.Email {
					t.Errorf("Expected email %v, got %v", tt.contact.Email, contact.Email)
				}
				if contact.Mobile != tt.contact.Mobile {
					t.Errorf("Expected mobile %v, got %v", tt.contact.Mobile, contact.Mobile)
				}
				if contact.FirstName != tt.contact.FirstName {
					t.Errorf("Expected first name %v, got %v", tt.contact.FirstName, contact.FirstName)
				}
				if contact.LastName != tt.contact.LastName {
					t.Errorf("Expected last name %v, got %v", tt.contact.LastName, contact.LastName)
				}
				if contact.CountryCode != tt.contact.CountryCode {
					t.Errorf("Expected country code %v, got %v", tt.contact.CountryCode, contact.CountryCode)
				}
			}
		})
	}
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
	testCases := []struct {
		name           string
		updatedContact models.Contact
		expectedError  bool
	}{
		{
			name: "ValidUpdate",
			updatedContact: models.Contact{
				UUID:        existingUUID,
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "+1122334455",
				Email:       "test.new@example.com",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			expectedError: false,
		},
		{
			name: "InvalidEmailFormat",
			updatedContact: models.Contact{
				UUID:        existingUUID,
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "+1122334455",
				Email:       "invalid-email",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			expectedError: true,
		},
		{
			name: "InvalidMobileFormat",
			updatedContact: models.Contact{
				UUID:        existingUUID,
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "1122334455",
				Email:       "valid@example.com",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			expectedError: true,
		},
		{
			name: "NonExistentUUID",

			updatedContact: models.Contact{
				UUID:        uuid.New(),
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "+1122334455",
				Email:       "test1.test1@example.com",
				CountryCode: "USA",
				ListID:      list.ID,
			},
			expectedError: true,
		},
		{
			name: "InvalidCountryCode",
			updatedContact: models.Contact{
				UUID:        existingUUID,
				FirstName:   "Test",
				LastName:    "Test",
				Mobile:      "+1122334455",
				Email:       "test1.test1@example.com",
				CountryCode: "US",
				ListID:      list.ID,
			},
			expectedError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			err := service.UpdateContact(tt.updatedContact)
			if tt.expectedError {
				if err == nil {
					t.Fatal("Expected error, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				var updatedContact models.Contact
				if err := db.Where("uuid = ?", tt.updatedContact.UUID).First(&updatedContact).Error; err != nil {
					t.Fatalf("Failed to find updated contact: %v", err)
				}
				if updatedContact.Email != tt.updatedContact.Email {
					t.Errorf("Expected email %v, got %v", tt.updatedContact.Email, updatedContact.Email)
				}
			}
		})
	}
}

func TestContactService_DeleteContact(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	service := services.NewContactService(repo)

	testLists := []models.List{
		{UUID: uuid.New(), Name: "Test List"},
	}
	for _, list := range testLists {
		if err := db.Create(&list).Error; err != nil {
			t.Fatalf("Could not create test list: %v", err)
		}
	}

	testContacts := []models.Contact{
		{
			UUID:        uuid.New(),
			FirstName:   "Test 1",
			LastName:    "Test",
			Mobile:      "+1234567890",
			Email:       "test1@example.com",
			CountryCode: "USA",
			ListID:      testLists[0].ID,
		},
		{
			UUID:        uuid.New(),
			FirstName:   "Test 2",
			LastName:    "Test",
			Mobile:      "+1987654321",
			Email:       "test2@example.com",
			CountryCode: "USA",
			ListID:      testLists[0].ID,
		},
	}

	for _, contact := range testContacts {
		if err := db.Create(&contact).Error; err != nil {
			t.Fatalf("Could not create test contact: %v", err)
		}
	}

	testCases := []struct {
		name        string
		contactUUID uuid.UUID
		expectError bool
		shouldExist bool
	}{
		{
			name:        "DeleteExistingContact",
			contactUUID: testContacts[0].UUID,
			expectError: false,
			shouldExist: false,
		},
		{
			name:        "DeleteAlreadyDeletedContact",
			contactUUID: testContacts[1].UUID,
			expectError: true,
			shouldExist: false,
		},
		{
			name:        "DeleteNonExistentContact",
			contactUUID: uuid.New(),
			expectError: true,
			shouldExist: false,
		},
	}

	service.DeleteContact(testContacts[1].UUID)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			err := service.DeleteContact(tt.contactUUID)
			if (err != nil) != tt.expectError {
				t.Fatalf("Expected error: %v, got %v", tt.expectError, err)
			}

			var deletedContact models.Contact
			result := db.Where("uuid = ?", tt.contactUUID).First(&deletedContact)
			if (result.Error == nil) != tt.shouldExist {
				t.Errorf("Expected contact existence: %v, but found %v", tt.shouldExist, result.Error == nil)
			}
		})
	}
}
