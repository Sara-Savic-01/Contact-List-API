package handlers
import (
	"testing"
	"contact-list-api-1/handlers"
	"contact-list-api-1/services"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"strings"
	"errors"
	"fmt"
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
		    Mobile:      "+0987654321",
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

	    t.Run("WithoutFilterAndPagination", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/contacts", nil)
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

		if len(gotContacts) != len(contacts) {
		    t.Errorf("Expected %d contacts, got %d", len(contacts), len(gotContacts))
		}
	    })

	    t.Run("WithNameFilter", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/contacts?name=Test", nil)
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

		if len(gotContacts) != 3 {
		    t.Errorf("Expected 3 contacts with name 'Test', got %d", len(gotContacts))
		}
	    })

	    t.Run("WithMobileFilter", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/contacts?mobile=%2B1234567890", nil) 
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

		if len(gotContacts) != 1 || gotContacts[0].Mobile != "+1234567890" {
		    t.Errorf("Expected 1 contact with mobile '+1234567890', got %v", gotContacts)
		}
	    })

	    t.Run("WithPagination", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/contacts?page=2&pageSize=1", nil)
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

		if len(gotContacts) != 1 || gotContacts[0].FirstName != "Test 1" { 
		    t.Errorf("Expected 1 contact on page 2 with first name 'Test 1', got %v", gotContacts)
		}
	    })

	    t.Run("WithNameFilterAndPagination", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/contacts?name=Test&page=1&pageSize=1", nil)
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

		if len(gotContacts) != 1 || gotContacts[0].FirstName != "Test" {
		    t.Errorf("Expected 1 contact on page 1 with first name 'Test', got %v", gotContacts)
		}
	    })
}

func TestGetContactByUUID(t *testing.T) {
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

		handler := handlers.NewContactHandler(service)

		req, err := http.NewRequest("GET", "/contacts/get/"+testUUID.String(), nil)
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.GetContactByUUID(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
		}

		var gotContact models.Contact
		if err := json.NewDecoder(rr.Body).Decode(&gotContact); err != nil {
			t.Fatalf("Could not decode response body: %v", err)
		}

		if gotContact.UUID != testContact.UUID || gotContact.Email != testContact.Email {
			t.Errorf("Expected contact %+v, got %+v", testContact, gotContact)
		}
		t.Run("NonExistentUUID", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/contacts/get/"+uuid.New().String(), nil)
			if err != nil {
			    t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.GetContactByUUID(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
			    t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
			}
	    	})

	    	t.Run("MissingUUID", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/contacts/get/", nil)
			if err != nil {
			    t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.GetContactByUUID(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
			    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
	    	})
		t.Run("InvalidUUID", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/contacts/get/invalid-uuid", nil)
			if err != nil {
			    t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.GetContactByUUID(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
			    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
    		})
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
		newContact := fmt.Sprintf(`{
		"first_name": "Test",
		"last_name": "Test",
		"mobile": "+123456789",
		"email": "test.test@example.com",
		"country_code": "USA",
		"list_id": %d
		}`, list.ID)

		req, err := http.NewRequest("POST", "/contacts/create", strings.NewReader(newContact))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.CreateContact(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
		}

		var createdContact models.Contact
		if err := json.NewDecoder(rr.Body).Decode(&createdContact); err != nil {
			t.Fatalf("Could not decode response body: %v", err)
		}

		if createdContact.FirstName != "Test" || createdContact.Email != "test.test@example.com" {
			t.Errorf("Expected contact 'Test Test', got %+v", createdContact)
		}
		t.Run("InvalidListID", func(t *testing.T) {
			invalidListID := uint(111111) 
			invalidContact := fmt.Sprintf(`{
			    "first_name": "Test",
			    "last_name": "Test",
			    "mobile": "+987654321",
			    "email": "test.test@example.com",
			    "country_code": "USA",
			    "list_id": %d
			}`, invalidListID)

			req, err := http.NewRequest("POST", "/contacts/create", strings.NewReader(invalidContact))
			if err != nil {
			    t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.CreateContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest{
			    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
    		})
		t.Run("InvalidJSONFormat", func(t *testing.T) {
			req, err := http.NewRequest("POST", "/contacts/create", strings.NewReader(`{invalid json}`))
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.CreateContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
		})

		t.Run("MissingRequiredFields", func(t *testing.T) {
			incompleteContact := fmt.Sprintf(`{
				"first_name": "Test"
			}`)

			req, err := http.NewRequest("POST", "/contacts/create", strings.NewReader(incompleteContact))
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.CreateContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
		})
		t.Run("EmptyAndInvalidRequiredFields", func(t *testing.T) {
			incompleteContact := fmt.Sprintf(`{
				"first_name": "",
				"last_name": "",
				"mobile": "",
				"email": "",
				"country_code": "US",
				
				}`)

			req, err := http.NewRequest("POST", "/contacts/create", strings.NewReader(incompleteContact))
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.CreateContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
		})
		
}
func TestUpdateContact(t *testing.T) {
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
	    testContact := models.Contact{
		FirstName:   "test",
		LastName:    "test",
		Mobile:      "+123456789",
		Email:       "test.test@example.com",
		CountryCode: "USA",
		ListID:      list.ID,
	    }
	    if err := db.Create(&testContact).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	    }

	    handler := handlers.NewContactHandler(service)

	    updatedContact := fmt.Sprintf(`{
	    "first_name": "test",
	    "last_name": "test",
	    "mobile": "+123456789",
	    "email": "test@example.com",
	    "country_code": "USA",
	    "list_id": %d
	    }`, list.ID)

	    req, err := http.NewRequest("PUT", "/contacts/update/"+testContact.UUID.String(), strings.NewReader(updatedContact))
	    if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	    }

	    rr := httptest.NewRecorder()
	    handler.UpdateContact(rr, req)

	    if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, status)
	    }

	    var updated models.Contact
	    if err := db.First(&updated, "email = ?", "test@example.com").Error; err != nil {
		t.Fatalf("Could not retrieve updated contact: %v", err)
	    }

	    if updated.FirstName != "test" || updated.Mobile != "+123456789" {
		t.Errorf("Expected contact 'test test', got %+v", updated)
	    }
	    t.Run("InvalidListID", func(t *testing.T) {
		invalidListID := uint(999999) 
		updatedContact := fmt.Sprintf(`{
		    "first_name": "test",
		    "last_name": "test",
		    "mobile": "+123456789",
		    "email": "test@example.com",
		    "country_code": "USA",
		    "list_id": %d
		}`, invalidListID)

		req, err := http.NewRequest("PUT", "/contacts/update/"+testContact.UUID.String(), strings.NewReader(updatedContact))
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.UpdateContact(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
		    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
    	    })
	    t.Run("InvalidJSONFormat", func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/contacts/update/"+testContact.UUID.String(), strings.NewReader(`{invalid json}`))
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.UpdateContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
	    })
	    
	     t.Run("NonExistentUUID", func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/contacts/update/"+uuid.New().String(), strings.NewReader(updatedContact))
			if err != nil {
			    t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.UpdateContact(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
			    t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
			}
	    })

	    t.Run("MissingUUID", func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/contacts/update/", strings.NewReader(updatedContact))
			if err != nil {
			    t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.UpdateContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
			    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
	    })
	    t.Run("InvalidUUID", func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/contacts/update/invalid-uuid", strings.NewReader(updatedContact))
			if err != nil {
			    t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.UpdateContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
			    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
    	    })
			    

      
}
func TestDeleteContact(t *testing.T) {
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
		testUUID := uuid.New()
		testContact := models.Contact{
			UUID:        testUUID,
			FirstName:   "test",
			LastName:    "test",
			Mobile:      "+123456789",
			Email:       "test.test@example.com",
			CountryCode: "USA",
			ListID:      list.ID,
		}
		if err := db.Create(&testContact).Error; err != nil {
			t.Fatalf("Could not create test data: %v", err)
		}

		handler := handlers.NewContactHandler(service)

		req, err := http.NewRequest("DELETE", "/contacts/delete/"+testUUID.String(), nil)
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.DeleteContact(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, status)
		}

		var deletedContact models.Contact
		if err := db.First(&deletedContact, "uuid = ?", testUUID).Error; err == nil {
			t.Errorf("Expected contact to be deleted, but it still exists in the database: %+v", deletedContact)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("Unexpected error while checking if contact was deleted: %v", err)
		}
		t.Run("UUIDNotFound", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/contacts/delete/"+uuid.New().String(), nil)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.DeleteContact(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
			}
		})
		t.Run("InvalidUUIDFormat", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/contacts/delete/invalid-uuid", nil)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.DeleteContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
		})
		t.Run("MissingUUID", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/contacts/delete/", nil)
			if err != nil {
			    t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.DeleteContact(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
			    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
			}
	    	})
}
