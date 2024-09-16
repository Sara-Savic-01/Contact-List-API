package handlers

import (
	"contact-list-api-1/handlers"
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"contact-list-api-1/services"
	"contact-list-api-1/tests"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setTestDB(t *testing.T) (*gorm.DB, func()) {
	db := tests.SetupTestDB(t)
	cleanup := func() {
		tests.TearDownTestDB(t, db)
	}
	return db, cleanup
}
func TestGetAllLists(t *testing.T) {
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

	testCases := []struct {
		name               string
		query              string
		expectedCount      int
		expectedListName   string
		expectedStatusCode int
	}{
		{
			name:               "WithoutFilterAndPagination",
			query:              "",
			expectedCount:      len(lists),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "WithNameFilter",
			query:              "?name=Another",
			expectedCount:      1,
			expectedListName:   "Another List",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "WithPagination",
			query:              "?page=2&pageSize=1",
			expectedCount:      1,
			expectedListName:   "List 2",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "WithNameFilterAndPagination",
			query:              "?name=List&page=1&pageSize=2",
			expectedCount:      2,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/lists"+tt.query, nil)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.GetAllLists(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}

			var gotLists []models.List
			if err := json.NewDecoder(rr.Body).Decode(&gotLists); err != nil {
				t.Fatalf("Could not decode response body: %v", err)
			}

			if len(gotLists) != tt.expectedCount {
				t.Errorf("Expected %d lists, got %d", tt.expectedCount, len(gotLists))
			}
			if tt.expectedListName != "" && (len(gotLists) > 0 && gotLists[0].Name != tt.expectedListName) {
				t.Errorf("Expected list name '%s', got '%s'", tt.expectedListName, gotLists[0].Name)
			}
		})
	}
}
func TestGetListByUUID(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)
	handler := handlers.NewListHandler(service)

	testUUID := uuid.New()
	testList := models.List{
		UUID: testUUID,
		Name: "Test List",
	}
	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	testCases := []struct {
		name               string
		uuid               string
		expectedStatusCode int
		expectedList       *models.List
	}{
		{
			name:               "ValidUUID",
			uuid:               testUUID.String(),
			expectedStatusCode: http.StatusOK,
			expectedList:       &testList,
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
			handler.GetListByUUID(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}

			if tt.expectedList != nil {
				var gotList models.List
				if err := json.NewDecoder(rr.Body).Decode(&gotList); err != nil {
					t.Fatalf("Could not decode response body: %v", err)

				}
				if !reflect.DeepEqual(gotList, *tt.expectedList) {
					t.Errorf("Expected list %+v, got %+v", *tt.expectedList, gotList)
				}
			}
		})
	}
}
func TestCreateList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)
	handler := handlers.NewListHandler(service)

	testCases := []struct {
		name               string
		body               string
		expectedStatusCode int
		expectedListName   string
	}{
		{
			name:               "ValidCreate",
			body:               `{"name": "New Test List"}`,
			expectedStatusCode: http.StatusCreated,
			expectedListName:   "New Test List",
		},
		{
			name:               "InvalidJSONFormat",
			body:               `{"name": "Invalid List"`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "MissingNameField",
			body:               `{}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "EmptyNameField",
			body:               `{"name":""}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/lists", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.CreateList(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}

			if tt.expectedStatusCode == http.StatusCreated {
				var createdList models.List
				if err := json.NewDecoder(rr.Body).Decode(&createdList); err != nil {
					t.Fatalf("Could not decode response body: %v", err)
				}
				if createdList.Name != tt.expectedListName {
					t.Errorf("Expected list name '%s', got '%s'", tt.expectedListName, createdList.Name)
				}
			}
		})
	}
}

func TestUpdateList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)
	handler := handlers.NewListHandler(service)

	testList := models.List{
		//UUID: uuid.New(),
		Name: "Old Test List",
	}

	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}
	updatedList := `{"name": "Updated Test List"}`
	testCases := []struct {
		name               string
		uuid               string
		body               string
		expectedStatusCode int
	}{
		{
			name:               "ValidUpdate",
			uuid:               testList.UUID.String(),
			body:               updatedList,
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "InvalidJSONFormat",
			uuid:               testList.UUID.String(),
			body:               `{"name": "Invalid List"`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "InvalidUUID",
			uuid:               "invalid-uuid",
			body:               updatedList,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "NonExistentUUID",
			uuid:               uuid.New().String(),
			body:               updatedList,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "MissingUUID",
			uuid:               "",
			body:               updatedList,
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

			handler.UpdateList(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}
		})
	}
}
func TestDeleteList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)
	handler := handlers.NewListHandler(service)

	testList := models.List{
		Name: "Test List to Delete",
	}

	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	testCases := []struct {
		name               string
		uuid               string
		expectedStatusCode int
	}{
		{
			name:               "ValidDelete",
			uuid:               testList.UUID.String(),
			expectedStatusCode: http.StatusNoContent,
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
			req, err := http.NewRequest("DELETE", "/", nil)
			req.SetPathValue("uuid", tt.uuid)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.DeleteList(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, status)
			}
		})
	}
}
