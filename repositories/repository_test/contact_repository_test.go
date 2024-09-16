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
		{UUID: uuid.New(), FirstName: "Test 1", LastName: "Contact 1", Mobile: "+1987654321", Email: "test.test@example.com", CountryCode: "USA", ListID: list.ID},
	}
	for _, contact := range testContacts {
		if err := db.Create(&contact).Error; err != nil {
			t.Fatalf("Could not create test contact: %v", err)
		}
	}

	testCases := []struct {
		name          string
		filterName    string
		filterEmail   string
		filterMobile  string
		limit         int
		offset        int
		expectedCount int
		expectedName  string
	}{
		{
			name:          "WithouFilterAndPagination",
			filterName:    "",
			filterEmail:   "",
			filterMobile:  "",
			limit:         0,
			offset:        0,
			expectedCount: len(testContacts),
		},
		{
			name:          "WithFilter",
			filterName:    "Contact",
			filterMobile:  "",
			filterEmail:   "",
			limit:         0,
			offset:        0,
			expectedCount: 2,
			expectedName:  "Contact",
		},
		{
			name:          "WithPagination",
			filterName:    "",
			filterMobile:  "",
			filterEmail:   "",
			limit:         0,
			offset:        0,
			expectedCount: 2,
		},
		{
			name:          "WithPagination",
			filterName:    "",
			filterMobile:  "",
			filterEmail:   "",
			limit:         1,
			offset:        1,
			expectedCount: 1,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			contacts, err := repo.GetAll(tt.filterName, tt.filterMobile, tt.filterEmail, tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)

			}
			if len(contacts) != tt.expectedCount {
				t.Errorf("Expected %d contacts, got %d", tt.expectedCount, len(contacts))
			}
			if tt.expectedName != "" && len(contacts) > 0 && contacts[0].LastName != tt.expectedName {
				t.Errorf("Expected contact with first name '%s', got %s", tt.expectedName, contacts[0].LastName)
			}
		})
	}
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

	testCases := []struct {
		name          string
		contact       models.Contact
		expectedError bool
	}{
		{
			name:          "CreateValidContact",
			contact:       testContact,
			expectedError: false,
		},
		{
			name:          "DuplicateContact",
			contact:       testContact,
			expectedError: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.Create(tt.contact); (err != nil) != tt.expectedError {
				t.Fatalf("Expected error:%v, got %v", tt.expectedError, err)

			}
			if !tt.expectedError {
				var createdContact models.Contact

				if err := db.Where("uuid=?", testUUID).First(&createdContact).Error; err != nil {
					t.Fatalf("Test data not found in database:%v", err)
				}
			}
		})
	}
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
		Mobile:      "+1987654321",
		Email:       "test.contact@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
	}

	if err := db.Create(&testContact).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	testCases := []struct {
		name          string
		uuid          uuid.UUID
		expectedError bool
		expectedUUID  uuid.UUID
	}{
		{
			name:          "GetExistingContact",
			uuid:          testUUID,
			expectedError: false,
			expectedUUID:  testUUID,
		},
		{
			name:          "NonExistentUUID",
			uuid:          uuid.New(),
			expectedError: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			contact, err := repo.GetByUUID(tt.uuid)
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error:%v, got %v", tt.expectedError, err)

			}
			if !tt.expectedError && (contact == nil || contact.UUID != tt.expectedUUID) {
				t.Errorf("Expected contact with UUID %v, got %v", tt.expectedUUID, contact)
			}
		})
	}
}

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

	testCases := []struct {
		name             string
		contact          models.Contact
		expectedError    bool
		expectedLastName string
	}{
		{
			name:             "UpdateExistingContact",
			contact:          models.Contact{UUID: testUUID, LastName: "Contact 1"},
			expectedError:    false,
			expectedLastName: "Contact 1",
		},
		{
			name:          "UpdateNonExistentContact",
			contact:       models.Contact{UUID: uuid.New(), FirstName: "Non", LastName: "Existent"},
			expectedError: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.contact)
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error:%v, got %v", tt.expectedError, err)

			}
			if !tt.expectedError {
				var updatedContact models.Contact
				if err := db.Where("uuid = ?", testUUID).First(&updatedContact).Error; err != nil {
					t.Fatalf("Contact not found after update: %v", err)
				}
				if updatedContact.LastName != tt.expectedLastName {
					t.Errorf("Expected last name to be '%s', got %s", tt.expectedLastName, updatedContact.LastName)
				}
			}
		})
	}

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

	tests := []struct {
		name        string
		UUID        uuid.UUID
		expectError bool
	}{
		{
			name:        "DeleteExistingContact",
			UUID:        testUUID,
			expectError: false,
		},
		{
			name:        "DeleteNonExistentContact",
			UUID:        uuid.New(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.UUID)
			if (err != nil) != tt.expectError {
				t.Fatalf("Expected error: %v, got %v", tt.expectError, err)
			}
			if !tt.expectError {
				var deletedContact models.Contact
				result := db.Where("uuid = ?", tt.UUID).First(&deletedContact)
				if result.Error == nil {
					t.Errorf("Expected record to be deleted, but it still exists")
				}
			}
		})
	}

}
