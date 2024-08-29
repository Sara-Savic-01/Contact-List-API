package config

import (
    
    "os"
    "testing"
    "contact-list-api-1/config"
)

func TestLoadConfig_ValidFile(t *testing.T) {

    	testConfig := `{
		"db": {
		    "user": "admin",
		    "password": "secret",
		    "host": "localhost",
		    "name": "test_db"
		},
		"auth_token": "my-secret-token"
	    }`
    	file, err := os.Create("test_config.json")
    	if err != nil {
        	t.Fatalf("Failed to create test config file: %v", err)
    	}
    	defer os.Remove(file.Name())
    	file.WriteString(testConfig)
    	file.Close()

    
    	config, err := config.LoadConfig(file.Name())
    	if err != nil {
        	t.Fatalf("Failed to load config: %v", err)
    	}

    	if config.DB.User != "admin" {
        	t.Errorf("Expected DB user 'admin', got %s", config.DB.User)
    	}
    	if config.DB.Password != "secret" {
        	t.Errorf("Expected DB password 'secret', got %s", config.DB.Password)
    	}
    	if config.DB.Host != "localhost" {
        	t.Errorf("Expected DB host 'localhost', got %s", config.DB.Host)
    	}
    	if config.DB.Name != "test_db" {
        	t.Errorf("Expected DB name 'test_db', got %s", config.DB.Name)
    	}
    	if config.AuthToken != "my-secret-token" {
        	t.Errorf("Expected auth token 'my-secret-token', got %s", config.AuthToken)
    	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {
    
    	invalidConfig := `{"db": {"user": "admin", "password": "secret", "host": "localhost", "name": 1234}}`
    	file, err := os.Create("invalid_config.json")
    	if err != nil {
        	t.Fatalf("Failed to create invalid config file: %v", err)
    	}
    	defer os.Remove(file.Name())
    	ile.WriteString(invalidConfig)
    	file.Close()

    	_, err = config.LoadConfig(file.Name())
    	if err == nil {
        	t.Fatal("Expected error when loading from invalid config file, got nil")
    	}
}

func TestLoadConfig_NoFile(t *testing.T) {
    
    	_, err := config.LoadConfig("no_file.json")
    	if err == nil {
        	t.Fatal("Expected error when loading from non-existent config file, got nil")
    	}
}

func TestLoadTestConfig_ValidFile(t *testing.T) {
    
	testConfig := `{
		"db": {
		    "user": "test_user",
		    "password": "test_password",
		    "host": "test_host",
		    "name": "test_db"
		}
	}`
	file, err := os.Create("test_config1.json")
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(file.Name())
	file.WriteString(testConfig)
	file.Close()

	    
	config := config.LoadTestConfig(file.Name())
	if config.DB.User != "test_user" {
		t.Errorf("Expected DB user 'test_user', got %s", config.DB.User)
	}
	if config.DB.Password != "test_password" {
		t.Errorf("Expected DB password 'test_password', got %s", config.DB.Password)
	}
	if config.DB.Host != "test_host" {
		t.Errorf("Expected DB host 'test_host', got %s", config.DB.Host)
	}
	if config.DB.Name != "test_db" {
		t.Errorf("Expected DB name 'test_db', got %s", config.DB.Name)
	}
}



func TestLoadTestConfig_InvalidJSON(t *testing.T) {

	invalidJSON := `{
		"db": {
			"user": "testuser",
			"password": "testpass",
			"host": "",
			"name": "testname"
		}
	}`

	tmpFile, err := os.Create("invalid_config.json")
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	defer os.Remove("invalid_config.json") 

	_, err = tmpFile.WriteString(invalidJSON)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON to temporary config file: %v", err)
	}
	tmpFile.Close()

	config := config.LoadTestConfig("invalid_config.json")
	
	if config.DB.Host!=""{
		t.Errorf("Expected empty hostname, got %v", config.DB.Host)
	}
}
