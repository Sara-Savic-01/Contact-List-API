package services
import (
	"testing"
	"contact-list-api-1/services"
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

	t.Run("WithoutFilterAndPagination", func(t *testing.T) {
		results, err := service.GetAllLists("", 1, 10) 
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(results) != len(lists) {
		    t.Errorf("Expected %d lists, got %d", len(lists), len(results))
		}
	    })

	t.Run("WithFilter", func(t *testing.T) {
		results, err := service.GetAllLists("Family", 1, 10)
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(results) != 1 {
		    t.Errorf("Expected 1 list with name 'Family', got %d", len(results))
		}
		if results[0].Name != "Family" {
		    t.Errorf("Expected list with name 'Family', got %s", results[0].Name)
		}
	})

	t.Run("WithPagination", func(t *testing.T) {
		results, err := service.GetAllLists("", 2, 1) 
		if err != nil {
		    t.Fatalf("Expected no error, got %v", err)
		}
		if len(results) != 1 {
		    t.Errorf("Expected 1 list on page 2 with page size 1, got %d", len(results))
		}
		
	})
	
	

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

	result, err := service.GetListByUUID(testUUID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.UUID != testUUID {
		t.Errorf("Expected UUID %v, got %v", testUUID, result.UUID)
	}
	
	t.Run("NonExistentUUID", func(t *testing.T) {
		nonExistentUUID := uuid.New()
		result, err := service.GetListByUUID(nonExistentUUID)
		if err == nil {
		    t.Fatalf("Expected an error, got %v", err)
		}
		if result != nil {
		    t.Errorf("Expected no list, got %v", result)
		}
    })

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

    	err := service.CreateList(newList)
    	if err != nil {
        	t.Fatalf("Expected no error, got %v", err)
    	}

    	var list models.List
    	db.Where("uuid = ?", newList.UUID).First(&list)
    	if list.Name != newList.Name {
        	t.Errorf("Expected name %v, got %v", newList.Name, list.Name)
    	}
	t.Run("EmptyName", func(t *testing.T) {
    		newList := models.List{
        		UUID: uuid.New(),
        		Name: "", 
    		}

    		err := service.CreateList(newList)
    		if err == nil {
        		t.Fatalf("Expected an error for invalid name, got nil")
    		}
	})
	

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

    	list.Name = "Updated Name"
    	err := service.UpdateList(list)
    	if err != nil {
        	t.Fatalf("Expected no error, got %v", err)
    	}

    	var updatedList models.List
    	db.Where("uuid = ?", existingUUID).First(&updatedList)
    	if updatedList.Name != list.Name {
        	t.Errorf("Expected name %v, got %v", list.Name, updatedList.Name)
    	}
	t.Run("NonExistentUUID", func(t *testing.T) {
    		invalidUUID := uuid.New()
    		list := models.List{
        		UUID: invalidUUID,
        		Name: "Update Attempt",
    		}

    		err := service.UpdateList(list)
    		if err == nil {
        		t.Fatalf("Expected an error for invalid UUID, got nil")
    		}
	})
	
}

func TestListService_DeleteList(t *testing.T) {
    	db, cleanup := setTestDB(t)
    	defer cleanup()

    	repo := repositories.NewListRepository(db)
    	service := services.NewListService(repo)
	

    	testUUID := uuid.New()
    	list := models.List{
        	UUID: testUUID,
        	Name: "To be deleted",
    	}
    	db.Create(&list)

    	err := service.DeleteList(testUUID)
    	if err != nil {
        	t.Fatalf("Expected no error, got %v", err)
    	}

    	var deletedList models.List
    	if result := db.Where("uuid = ?", testUUID).First(&deletedList); result.Error == nil {
        	t.Errorf("Expected record to be deleted, but it still exists")
    	}
	
	t.Run("DeleteAlreadyDeleted", func(t *testing.T) {
		    validUUID := uuid.New()
		    list := models.List{
			UUID: validUUID,
			Name: "To be deleted",
		    }
		    db.Create(&list)
		    db.Delete(&list) 

		    err := service.DeleteList(validUUID)
		    if err == nil {
			t.Fatalf("Expected error when deleting already deleted list, got %v", err)
		    }

		    var deletedList models.List
		    if result := db.Where("uuid = ?", validUUID).First(&deletedList); result.Error == nil {
			t.Errorf("Expected record to be deleted, but it still exists")
		    }
	})
	t.Run("NonExistentList", func(t *testing.T) {
		nonExistentUUID := uuid.New()
		err := service.DeleteList(nonExistentUUID)
		if err == nil {
		    t.Fatalf("Expected an error, got %v", err)
		}
		
		var list models.List
		if result := db.Where("uuid = ?", nonExistentUUID).First(&list); result.Error == nil {
		    t.Errorf("Expected list to not exist, but it was found")
		}
    })
}