package repositories

import (
	"testing"

	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"contact-list-api-1/tests"

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

func TestListRepository_GetAll(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	lists := []models.List{
		{UUID: uuid.New(), Name: "List1"},
		{UUID: uuid.New(), Name: "List2"},
		{UUID: uuid.New(), Name: "List3"},
	}
	for i, list := range lists {
		if err := db.Create(&list).Error; err != nil {
			t.Fatalf("Could not create list: %v", err)
		}
		lists[i] = list
	}
	testCases := []struct {
		name          string
		filter        string
		limit         int
		offset        int
		expectedCount int
		expectedName  string
	}{
		{
			name:          "WithoutFilterAndPagination",
			filter:        "",
			limit:         0,
			offset:        0,
			expectedCount: 3,
		},
		{
			name:          "WithNameFilter",
			filter:        "List1",
			limit:         0,
			offset:        0,
			expectedCount: 1,
			expectedName:  "List1",
		},
		{
			name:          "WithPagination",
			filter:        "",
			limit:         2,
			offset:        1,
			expectedCount: 2,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			lists, err := repo.GetAll(tt.filter, tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)

			}
			if len(lists) != tt.expectedCount {
				t.Errorf("Expected %d lists, got %d", tt.expectedCount, len(lists))
			}
			if tt.expectedName != "" && len(lists) > 0 && lists[0].Name != tt.expectedName {
				t.Errorf("Expected list with name '%s', got '%s'", tt.expectedName, lists[0].Name)
			}
		})
	}
}

func TestListRepository_Create(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	testUUID := uuid.New()
	testList := models.List{UUID: testUUID, Name: "Test List"}
	testCases := []struct {
		name          string
		list          models.List
		expectedError bool
		expectedName  string
	}{
		{
			name:          "CreateValidList",
			list:          testList,
			expectedError: false,
			expectedName:  "Test List",
		},
		{
			name:          "DuplicateList",
			list:          testList,
			expectedError: true,
			expectedName:  "Test List",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.Create(tt.list); (err != nil) != tt.expectedError {
				t.Fatalf("Expected error:%v, got %v", tt.expectedError, err)

			}
			if !tt.expectedError {
				var fetchedList models.List
				db.Where("uuid=?", tt.list.UUID).First(&fetchedList)
				if fetchedList.Name != tt.expectedName {
					t.Errorf("Expected list name to be '%s', got %s", tt.expectedName, fetchedList.Name)
				}
			}
		})
	}

}

func TestListRepository_GetByUUID(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)

	testUUID := uuid.New()
	testList := models.List{
		UUID: testUUID,
		Name: "Test List",
	}

	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	testCases := []struct {
		name          string
		uuid          uuid.UUID
		expectedError bool
		expectedUUID  uuid.UUID
	}{
		{
			name:          "GetExistingUUID",
			uuid:          testUUID,
			expectedError: false,
			expectedUUID:  testUUID,
		},
		{
			name:          "GetNonExistentUUID",
			uuid:          uuid.New(),
			expectedError: true,
			expectedUUID:  testUUID,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			list, err := repo.GetByUUID(tt.uuid)
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error:%v, got %v", tt.expectedError, err)

			}
			if !tt.expectedError && (list == nil || list.UUID != tt.expectedUUID) {
				t.Errorf("Expected list with UUID %v, got %v", tt.expectedUUID, list.UUID)
			}
		})
	}
}

func TestListRepository_Update(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)

	testUUID := uuid.New()
	testList := models.List{
		UUID: testUUID,
		Name: "Original Name",
	}

	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	testCases := []struct {
		name          string
		list          models.List
		expectedError bool
		expectedName  string
	}{
		{
			name:          "UpdateExistingList",
			list:          models.List{UUID: testUUID, Name: "Updated Name"},
			expectedError: false,
			expectedName:  "Updated Name",
		},
		{
			name:          "UpdateNonExistentList",
			list:          models.List{UUID: uuid.New(), Name: "Non"},
			expectedError: true,
			expectedName:  "",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.list)
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error:%v, got %v", tt.expectedError, err)

			}
			if !tt.expectedError {
				var updatedList models.List
				db.Where("uuid=?", tt.list.UUID).First(&updatedList)
				if updatedList.Name != tt.expectedName {
					t.Errorf("Expected list name to be '%s', got %s", tt.expectedName, updatedList.Name)
				}
			}
		})
	}

}

func TestListRepository_Delete(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)

	testUUID := uuid.New()
	testList := models.List{
		UUID: testUUID,
		Name: "To be deleted",
	}

	if err := db.Create(&testList).Error; err != nil {
		t.Fatalf("Could not create test data: %v", err)
	}

	testCases := []struct {
		name          string
		uuid          uuid.UUID
		expectedError bool
	}{
		{
			name:          "DeleteExistingList",
			uuid:          testUUID,
			expectedError: false,
		},
		{
			name:          "DeleteNonexistentUUID",
			uuid:          uuid.New(),
			expectedError: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.uuid)
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error:%v, got %v", tt.expectedError, err)

			}
			if !tt.expectedError {
				var list models.List
				if result := db.Where("uuid=?", tt.uuid).First(&list); result.Error == nil {
					t.Errorf("Expected record to be deleted, but it still exists")
				}
			}
		})
	}
}
