package repositories
import (
	"testing"
	
	"contact-list-api-1/models"
	"contact-list-api-1/repositories"
	"github.com/google/uuid"
	"contact-list-api-1/tests"
	"gorm.io/gorm"
)

func setTestDB(t *testing.T) (*gorm.DB, func()){
	db:=tests.SetupTestDB(t)
	cleanup:=func(){
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

	t.Run("NoFilterNoPagination", func(t *testing.T) {
		lists, err := repo.GetAll("", 0, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(lists) != 3 {
			t.Errorf("Expected 3 lists, got %d", len(lists))
		}
	})

	t.Run("WithNameFilter", func(t *testing.T) {
		lists, err := repo.GetAll("List1", 0, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(lists) != 1 {
			t.Errorf("Expected 1 list with name 'List1', got %d", len(lists))
		}
		if lists[0].Name != "List1" {
			t.Errorf("Expected list with name 'List1', got %s", lists[0].Name)
		}
	})

	t.Run("WithPagination", func(t *testing.T) {
		lists, err := repo.GetAll("", 2, 1) 
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(lists) != 2 {
			t.Errorf("Expected 2 lists, got %d", len(lists))
		}
		
	})
	t.Run("FilterNoMatch", func(t *testing.T) {
    	        lists, err := repo.GetAll("NonExistent", 0, 0)
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(lists) != 0 {
		    t.Errorf("Expected 0 lists, got %d", len(lists))
		}
	})
	t.Run("InvalidPagination", func(t *testing.T) {
	        lists, err := repo.GetAll("", -1, -1)
		if err == nil {
		    t.Fatalf("Expected an error, got none")
		}
		if len(lists) != 0 {
		    t.Errorf("Expected 0 lists, got %d", len(lists))
		}
	})

}

func TestListRepository_Create(t *testing.T) {
	db, cleanup := setTestDB(t)
	defer cleanup()

	repo := repositories.NewListRepository(db)
	
	
	testUUID := uuid.New()
	testList := models.List{
		UUID: testUUID,
		Name: "Test List",
	}

	err := repo.Create(testList)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var fetchedList models.List
	db.Where("uuid = ?", testUUID).First(&fetchedList)
	if fetchedList.Name != "Test List" {
		t.Errorf("Expected list name to be 'Test List', got %s", fetchedList.Name)
	}
	t.Run("MissingFields", func(t *testing.T) {
	    invalidList := testList
	    invalidList.Name = ""
	    err := repo.Create(invalidList)
	    if err == nil {
		t.Fatalf("Expected an error due to missing required fields, got none")
	    }
	})

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

	list, err := repo.GetByUUID(testUUID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if list.Name != "Test List" {
		t.Errorf("Expected list name to be 'Test List', got %s", list.Name)
	}
	t.Run("NonExistentUUID", func(t *testing.T) {
	    _, err := repo.GetByUUID(uuid.New())
	    if err == nil {
		t.Fatalf("Expected an error for non-existent UUID, got none")
	    }
	})

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

	testList.Name = "Updated Name"
	if err := repo.Update(testList); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var updatedList models.List
	db.Where("uuid = ?", testUUID).First(&updatedList)
	if updatedList.Name != "Updated Name" {
			t.Errorf("Expected list name to be 'Updated Name', got %s", updatedList.Name)
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

	if err := repo.Delete(testUUID); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var list models.List
	if result := db.Where("uuid = ?", testUUID).First(&list); result.Error == nil {
		t.Errorf("Expected record to be deleted, but it still exists")
	}
	t.Run("DeleteNonExistentUUID", func(t *testing.T) {
	    err := repo.Delete(uuid.New())
	    if err == nil {
		t.Fatalf("Expected an error for deleting non-existent UUID, got none")
	    }
	})
	
	
}
