package services

import (
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"contact-list-api-1/services"
	"contact-list-api-1/tests"
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
func TestListService_GetAllLists(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)

	lists := []models.List{
		{UUID: uuid.New(), Name: "Family"},
		{UUID: uuid.New(), Name: "Friends"},
		{UUID: uuid.New(), Name: "Work"},
	}
	for i, list := range lists {
		if err := db.Create(&list).Error; err != nil {
			t.Fatalf("Could not create test list: %v", err)
		}
		lists[i] = list
	}
	testCases := []struct {
		name          string
		filter        string
		page          int
		pageSize      int
		expectedCount int
		expectedName  string
	}{
		{
			name:          "WithoutFilterAndWithDefaultPagination",
			filter:        "",
			page:          1,
			pageSize:      10,
			expectedCount: len(lists),
		},
		{
			name:          "WithFilter",
			filter:        "Family",
			page:          1,
			pageSize:      10,
			expectedCount: 1,
			expectedName:  "Family",
		},
		{
			name:          "WithPagination",
			filter:        "",
			page:          2,
			pageSize:      1,
			expectedCount: 1,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.GetAllLists(tt.filter, tt.page, tt.pageSize)
			if err != nil {
				t.Fatalf("Excpected no error. got %v", err)
			}
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d lists, got %d", tt.expectedCount, len(results))
			}
			if tt.expectedName != "" && len(results) > 0 && results[0].Name != tt.expectedName {
				t.Errorf("Expected list with name '%s', got '%s'", tt.expectedName, results[0].Name)
			}
		})

	}

}

func TestListService_GetListByUUID(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)

	testUUID := uuid.New()
	list := models.List{
		UUID: testUUID,
		Name: "Test List",
	}
	db.Create(&list)

	testCases := []struct {
		name          string
		uuid          uuid.UUID
		expectedError bool
	}{
		{
			name:          "ExistingUUID",
			uuid:          testUUID,
			expectedError: false,
		},
		{
			name:          "NonExistentUUID",
			uuid:          uuid.New(),
			expectedError: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetListByUUID(tt.uuid)
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error:%v, got %v", tt.expectedError, err)
			}
			if !tt.expectedError {
				if result.UUID != tt.uuid {
					t.Errorf("Expected UUID %v, got %v", tt.uuid, result.UUID)
				}
			} else {
				if result != nil {
					t.Errorf("Expected no list, got %v", result)
				}
			}
		})
	}
}

func TestListService_CreateList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)

	newList := models.List{
		UUID: uuid.New(),
		Name: "New List",
	}
	testCases := []struct {
		name          string
		list          models.List
		expectedError bool
	}{
		{
			name:          "ValidList",
			list:          newList,
			expectedError: false,
		},
		{
			name: "EmptyName",
			list: models.List{
				UUID: uuid.New(),
				Name: "",
			},
			expectedError: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateList(tt.list)

			if tt.expectedError {
				if err == nil {
					t.Fatal("Expected error , got nil")
				}

			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				var createdList models.List
				db.Where("uuid = ?", tt.list.UUID).First(&createdList)

				if createdList.Name != tt.list.Name {
					t.Errorf("Expected name %v, got %v", tt.list.Name, createdList.Name)
				}

				if tt.list.UUID == uuid.Nil && createdList.UUID == uuid.Nil {
					t.Errorf("Expected a new UUID to be generated, got nil UUID")
				}
			}
		})
	}

}

func TestListService_UpdateList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)

	existingUUID := uuid.New()
	list := models.List{
		UUID: existingUUID,
		Name: "Original Name",
	}
	db.Create(&list)
	testCases := []struct {
		name          string
		updateList    models.List
		expectedError bool
	}{
		{
			name: "SuccessfulUpdate",

			updateList: models.List{
				UUID: existingUUID,
				Name: "Updated Name",
			},
			expectedError: false,
		},
		{
			name: "NonExistentUUID",

			updateList: models.List{
				UUID: uuid.New(),
				Name: "Update Attempt",
			},
			expectedError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			err := service.UpdateList(tt.updateList)

			if tt.expectedError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				var updatedList models.List
				db.Where("uuid = ?", tt.updateList.UUID).First(&updatedList)

				if updatedList.Name != tt.updateList.Name {
					t.Errorf("Expected name %v, got %v", tt.updateList.Name, updatedList.Name)
				}
			}
		})
	}

}

func TestListService_DeleteList(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	service := services.NewListService(repo)

	testLists := []models.List{
		{UUID: uuid.New(), Name: "To be deleted"},
		{UUID: uuid.New(), Name: "Existing List"},
	}

	for _, list := range testLists {
		if err := db.Create(&list).Error; err != nil {
			t.Fatalf("Could not create test list: %v", err)
		}
	}
	db.Delete(&testLists[0])
	testCases := []struct {
		name        string
		listUUID    uuid.UUID
		expectError bool
		shouldExist bool
	}{
		{
			name:        "DeleteExistingList",
			listUUID:    testLists[1].UUID,
			expectError: false,
			shouldExist: false,
		},
		{
			name:        "DeleteAlreadyDeletedList",
			listUUID:    testLists[0].UUID,
			expectError: true,
			shouldExist: false,
		},
		{
			name:        "DeleteNonExistentList",
			listUUID:    uuid.New(),
			expectError: true,
			shouldExist: false,
		},
	}
	service.DeleteList(testLists[0].UUID)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteList(tt.listUUID)
			if (err != nil) != tt.expectError {
				t.Fatalf("Expected error: %v, got %v", tt.expectError, err)
			}

			var deletedList models.List
			result := db.Where("uuid = ?", tt.listUUID).First(&deletedList)
			if (result.Error == nil) != tt.shouldExist {
				t.Errorf("Expected list existence: %v, but found %v", tt.shouldExist, result.Error == nil)
			}
		})
	}
}
