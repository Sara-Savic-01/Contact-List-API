package handlers
import (
	"testing"
	"contact-list-api-1/handlers"
	"contact-list-api-1/services"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"github.com/google/uuid"
	"contact-list-api-1/tests"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"strings"
	"errors"
	
)

func setTestDB(t *testing.T) (*gorm.DB, func()){
	db:=tests.SetupTestDB(t)
	cleanup:=func(){
		tests.TearDownTestDB(t, db)
	}
	return db, cleanup
}
func TestGetAllLists(t *testing.T){
	
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)
	handler := handlers.NewListHandler(service)

	lists := []models.List{
		{UUID: uuid.New(), Name: "List 1"},
		{UUID: uuid.New(), Name: "List 2"},
		{UUID: uuid.New(), Name: "Another List"},
	}
	for _, list := range lists {
		if err := db.Create(&list).Error; err != nil {
		    t.Fatalf("Could not create test list: %v", err)
		}
	}

	t.Run("WithoutFilterAndPagination", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lists", nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.GetAllLists(rr, req)

		if status := rr.Code; status != http.StatusOK {
		    t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
		}

		var gotLists []models.List
		if err := json.NewDecoder(rr.Body).Decode(&gotLists); err != nil {
		    t.Fatalf("Could not decode response body: %v", err)
		}

		if len(gotLists) != len(lists) {
		    t.Errorf("Expected %d lists, got %d", len(lists), len(gotLists))
		}
	})

	t.Run("WithNameFilter", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lists?name=Another", nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.GetAllLists(rr, req)

		if status := rr.Code; status != http.StatusOK {
		    t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
		}

		var gotLists []models.List
		if err := json.NewDecoder(rr.Body).Decode(&gotLists); err != nil {
		    t.Fatalf("Could not decode response body: %v", err)
		}

		if len(gotLists) != 1 || gotLists[0].Name != "Another List" {
		    t.Errorf("Expected 1 list named 'Another List', got %v", gotLists)
		}
	})

	t.Run("WithPagination", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lists?page=2&pageSize=1", nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.GetAllLists(rr, req)

		if status := rr.Code; status != http.StatusOK {
		    t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
		}

		var gotLists []models.List
		if err := json.NewDecoder(rr.Body).Decode(&gotLists); err != nil {
		    t.Fatalf("Could not decode response body: %v", err)
		}

		if len(gotLists) != 1 || gotLists[0].Name != "List 2" { 
		    t.Errorf("Expected 1 list on page 2 with name 'List 2', got %v", gotLists)
		}
	})

	t.Run("WithNameFilterAndPagination", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lists?name=List&page=1&pageSize=2", nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.GetAllLists(rr, req)

		if status := rr.Code; status != http.StatusOK {
		    t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
		}

		var gotLists []models.List
		if err := json.NewDecoder(rr.Body).Decode(&gotLists); err != nil {
		    t.Fatalf("Could not decode response body: %v", err)
		}

		if len(gotLists) != 2 {
		    t.Errorf("Expected 2 lists, got %d", len(gotLists))
		}
	})
	
   
	
}
func TestGetListByUUID(t *testing.T){
	db, cleanup:=setTestDB(t)
	defer cleanup()
	repo:=repositories.NewListRepository(db)
	service:=services.NewListService(repo)
	testUUID:=uuid.New()
	testList:=models.List{
		UUID: testUUID,
		Name:"Test List",
	}
	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	
	handler := handlers.NewListHandler(service)

	
	req, err := http.NewRequest("GET", "/lists/get/"+testUUID.String(), nil)
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler.GetListByUUID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var gotList models.List
	if err := json.NewDecoder(rr.Body).Decode(&gotList); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}

	if gotList.UUID != testList.UUID || gotList.Name != testList.Name{
		t.Errorf("Expected list %+v, got %+v", testList, gotList)
	}
	t.Run("InvalidUUID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lists/get/invalid-uuid", nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.GetListByUUID(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
		    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
    	})

    	t.Run("NonExistentUUID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lists/get/"+uuid.New().String(), nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.GetListByUUID(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
		    t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
		}
    	})

    	t.Run("MissingUUID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lists/get/", nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.GetListByUUID(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
		    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
    	})
	
	
}
func TestCreateList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)
	
	handler := handlers.NewListHandler(service)


	newList := `{"name": "New Test List"}`

	req, err := http.NewRequest("POST", "/lists/create", strings.NewReader(newList))
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler.CreateList(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	var createdList models.List
	if err := json.NewDecoder(rr.Body).Decode(&createdList); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}

	if createdList.Name != "New Test List" {
		t.Errorf("Expected list name 'New Test List', got %+v", createdList)
	}
	t.Run("InvalidJSONFormat", func(t *testing.T) {
		invalidJSON := `{"name": "Invalid List"`

		req, err := http.NewRequest("POST", "/lists/create", strings.NewReader(invalidJSON))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.CreateList(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
	})

	t.Run("MissingNameField", func(t *testing.T) {
		invalidJSON := `{}`

		req, err := http.NewRequest("POST", "/lists/create", strings.NewReader(invalidJSON))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.CreateList(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
	})
	t.Run("EmptyNameField", func(t *testing.T) {
		invalidJSON := `{"name":""}`

		req, err := http.NewRequest("POST", "/lists/create", strings.NewReader(invalidJSON))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.CreateList(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
	})
 	
}

func TestUpdateList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)

	testList := models.List{
		Name:  "Old Test List",
	}

	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	handler := handlers.NewListHandler(service)

	updatedList := `{"name": "Updated Test List"}`

	req, err := http.NewRequest("PUT", "/lists/update/"+testList.UUID.String(), strings.NewReader(updatedList))
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler.UpdateList(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, status)
	}
	
	t.Run("InvalidJSONFormat", func(t *testing.T) {
		invalidJSON := `{"name": "Invalid List"`

		req, err := http.NewRequest("PUT", "/lists/update/"+testList.UUID.String(), strings.NewReader(invalidJSON))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.UpdateList(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
	})
	t.Run("InvalidUUID", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/lists/update/invalid-uuid", strings.NewReader(updatedList))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateList(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
	})

	t.Run("NonExistentUUID", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/lists/update/"+uuid.New().String(), strings.NewReader(updatedList))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateList(rr, req)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
		}
	})

	t.Run("MissingUUID", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/lists/update/", strings.NewReader(updatedList))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateList(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
	})
}
func TestDeleteList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)

	testUUID := uuid.New()
	testList := models.List{
		UUID:  testUUID,
		Name:  "Test List",
	}

	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}	
	handler := handlers.NewListHandler(service)

	req, err := http.NewRequest("DELETE", "/lists/delete/"+testUUID.String(), nil)
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler.DeleteList(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, status)
	}

	var deletedList models.List
	if err := db.First(&deletedList, "uuid = ?", testUUID).Error; err == nil {
		t.Errorf("Expected list to be deleted, but it still exists in the database: %+v", deletedList)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Unexpected error while checking if list was deleted: %v", err)
	}
	t.Run("UUIDNotFound", func(t *testing.T) {
		nonExistentUUID := uuid.New()

		req, err := http.NewRequest("DELETE", "/lists/delete/"+nonExistentUUID.String(), nil)
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.DeleteList(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
		}
	})
	t.Run("InvalidUUIDFormat", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/lists/delete/invalid-uuid", nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.DeleteList(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
		    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
    	})
	t.Run("MissingUUID", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/lists/delete/", nil)
		if err != nil {
		    t.Fatalf("Could not create HTTP request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.DeleteList(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
		    t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
		}
    	})
}
	
	