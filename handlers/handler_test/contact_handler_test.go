package handlers

import (
	"contact-list-api-1/handlers"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"contact-list-api-1/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGetAllContacts(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	service := services.NewContactService(repo)
	handler := handlers.NewContactHandler(service)

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
			FirstName:   "Test",
			LastName:    "Test",
			Mobile:      "+1234567890",
			Email:       "test.test@example.com",
			CountryCode: "USA",
			ListID:      list.ID,
		},
		{
			UUID:        uuid.New(),
			FirstName:   "Test 1",
			LastName:    "Test 1",
			Mobile:      "+1987654321",
			Email:       "test1.test1@example.com",
			CountryCode: "USA",
			ListID:      list.ID,
		},
		{
			UUID:        uuid.New(),
			FirstName:   "Test 2",
			LastName:    "Test 2",
			Mobile:      "+1122334455",
			Email:       "test2.test2@example.com",
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
		name           string
		query          string
		expectedCount  int
		expectedName   string
		expectedMobile string
	}{
		{
			name:          "WithoutFilterAndPagination",
			query:         "",
			expectedCount: len(contacts),
		},
		{
			name:          "WithNameFilter",
			query:         "?name=Test",
			expectedCount: 3,
		},
		{
			name:           "WithMobileFilter",
			query:          "?mobile=%2B1234567890",
			expectedCount:  1,
			expectedMobile: "+1234567890",
		},
		{
			name:          "WithPagination",
			query:         "?page=2&pageSize=1",
			expectedCount: 1,
		},
		{
			name:          "WithNameFilterAndPagination",
			query:         "?name=Test&page=1&pageSize=1",
			expectedCount: 1,
			expectedName:  "Test",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/contacts"+tt.query, nil)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.GetAllContacts(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
			}

			var gotContacts []models.Contact
			if err := json.NewDecoder(rr.Body).Decode(&gotContacts); err != nil {
				t.Fatalf("Could not decode response body: %v", err)
			}

			if len(gotContacts) != tt.expectedCount {
				t.Errorf("Expected %d contacts, got %d", tt.expectedCount, len(gotContacts))
			}

			if tt.expectedMobile != "" && (len(gotContacts) > 0 && gotContacts[0].Mobile != tt.expectedMobile) {
				t.Errorf("Expected contact with mobile '%s', got %v", tt.expectedMobile, gotContacts)
			}

			if tt.expectedName != "" && (len(gotContacts) > 0 && gotContacts[0].FirstName != tt.expectedName) {
				t.Errorf("Expected contact with name '%s', got %v", tt.expectedName, gotContacts)
			}
		})
	}
}

func TestGetContactByUUID(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	service := services.NewContactService(repo)
	handler := handlers.NewContactHandler(service)
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
		LastName:    "Test",
		Mobile:      "+123456789",
		Email:       "test.test@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
	}
	if err := db.Create(&testContact).Error; err != nil {
		t.Fatalf("Could not create test contact: %v", err)
	}

	testCases := []struct {
		name               string
		uuid               string
		expectedStatusCode int
		expectedContact    *models.Contact
	}{
		{
			name:               "ValidUUID",
			uuid:               testUUID.String(),
			expectedStatusCode: http.StatusOK,
			expectedContact:    &testContact,
		},
		{
			name:               "InvalidUUID",
			uuid:               "invalid-uuid",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "NonExistentUUID",
			uuid:               uuid.New().String(),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "MissingUUID",
			uuid:               "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			req.SetPathValue("uuid", tt.uuid)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.GetContactByUUID(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}

			if tt.expectedContact != nil {
				var gotContact models.Contact
				if err := json.NewDecoder(rr.Body).Decode(&gotContact); err != nil {
					t.Fatalf("Could not decode response body: %v", err)
				}

				if !reflect.DeepEqual(gotContact, *tt.expectedContact) {
					t.Errorf("Expected contact %+v, got %+v", *tt.expectedContact, gotContact)
				}
			}
		})
	}
}

func TestCreateContact(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	service := services.NewContactService(repo)
	handler := handlers.NewContactHandler(service)

	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
	}
	if err := db.Create(&list).Error; err != nil {
		t.Fatalf("Could not create test list: %v", err)
	}

	testCases := []struct {
		name                string
		body                string
		expectedStatusCode  int
		expectedFirstName   string
		expectedLastName    string
		expectedMobile      string
		expectedEmail       string
		expectedCountryCode string
	}{
		{
			name: "ValidCreate",
			body: fmt.Sprintf(`{
				"first_name": "Test",
				"last_name": "Test",
				"mobile": "+123456789",
				"email": "test.test@example.com",
				"country_code": "USA",
				"list_id": %d
			}`, list.ID),
			expectedStatusCode:  http.StatusCreated,
			expectedFirstName:   "Test",
			expectedLastName:    "Test",
			expectedMobile:      "+123456789",
			expectedEmail:       "test.test@example.com",
			expectedCountryCode: "USA",
		},
		{
			name:               "InvalidJSONFormat",
			body:               `{"first_name": "Invalid Contact"`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "MissingRequiredFields",
			body:               `{"first_name": "Test"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "EmptyRequiredFields",
			body:               `{"first_name": "", "last_name": "", "mobile": "", "email": "", "country_code": ""}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/contacts", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.CreateContact(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}

			if tt.expectedStatusCode == http.StatusCreated {
				var createdContact models.Contact
				if err := json.NewDecoder(rr.Body).Decode(&createdContact); err != nil {
					t.Fatalf("Could not decode response body: %v", err)
				}
				if createdContact.FirstName != tt.expectedFirstName ||
					createdContact.LastName != tt.expectedLastName ||
					createdContact.Mobile != tt.expectedMobile ||
					createdContact.Email != tt.expectedEmail ||
					createdContact.CountryCode != tt.expectedCountryCode {
					t.Errorf("Expected contact with first name '%s', last name '%s', mobile '%s', email '%s', country code '%s', got first name '%s', last name '%s', mobile '%s', email '%s', country code '%s'",
						tt.expectedFirstName, tt.expectedLastName, tt.expectedMobile, tt.expectedEmail, tt.expectedCountryCode,
						createdContact.FirstName, createdContact.LastName, createdContact.Mobile, createdContact.Email, createdContact.CountryCode)
				}
			}
		})
	}
}
func TestUpdateContact(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	service := services.NewContactService(repo)
	handler := handlers.NewContactHandler(service)

	list := models.List{
		UUID: uuid.New(),
		Name: "Test List",
	}
	if err := db.Create(&list).Error; err != nil {
		t.Fatalf("Could not create test list: %v", err)
	}

	testContact := models.Contact{
		FirstName:   "test",
		LastName:    "test",
		Mobile:      "+123456789",
		Email:       "test.test@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
	}
	if err := db.Create(&testContact).Error; err != nil {
		t.Fatalf("Could not create test contact: %v", err)
	}

	updatedContact := fmt.Sprintf(`{
	    "first_name": "test",
	    "last_name": "test",
	    "mobile": "+123456789",
	    "email": "test.test1@example.com",
	    "country_code": "USA",
	    "list_id": %d
	}`, list.ID)

	testCases := []struct {
		name               string
		uuid               string
		body               string
		expectedStatusCode int
	}{
		{
			name:               "ValidUpdate",
			uuid:               testContact.UUID.String(),
			body:               updatedContact,
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "InvalidJSONFormat",
			uuid:               testContact.UUID.String(),
			body:               `{invalid json}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "NonExistentUUID",
			uuid:               uuid.New().String(),
			body:               updatedContact,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "InvalidUUIDFormat",
			uuid:               "invalid-uuid",
			body:               updatedContact,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "MissingUUID",
			uuid:               "",
			body:               updatedContact,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/", strings.NewReader(tt.body))
			req.SetPathValue("uuid", tt.uuid)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.UpdateContact(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}
		})
	}
}
func TestDeleteContact(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewContactRepository(db)
	service := services.NewContactService(repo)
	handler := handlers.NewContactHandler(service)

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
		LastName:    "Test",
		Mobile:      "+123456789",
		Email:       "test.test@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
	}
	if err := db.Create(&testContact).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	testCases := []struct {
		name               string
		uuid               string
		expectedStatusCode int
	}{
		{
			name:               "ValidDelete",
			uuid:               testUUID.String(),
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "UUIDNotFound",
			uuid:               uuid.New().String(),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "InvalidUUIDFormat",
			uuid:               "invalid-uuid",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "MissingUUID",
			uuid:               "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/", nil)
			req.SetPathValue("uuid", tt.uuid)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.DeleteContact(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}
		})
	}

}
