package config

import (
	"contact-list-api-1/config"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    *config.Config
		expectedErr bool
	}{
		{
			name: "Valid Config",
			fileContent: `{
				"db": {
					"user": "admin",
					"password": "secret",
					"host": "localhost",
					"name": "test_db"
				},
				"auth_token": "my-secret-token"
			}`,
			expected: &config.Config{
				DB: config.DBConfig{
					User:     "admin",
					Password: "secret",
					Host:     "localhost",
					Name:     "test_db",
				},
				AuthToken: "my-secret-token",
			},
			expectedErr: false,
		},
		{
			name: "Invalid JSON",
			fileContent: `{
				"db": {
					"user": "admin",
					"password": "secret",
					"host": "localhost",
					"name": 1234
				}}`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Missing File",
			fileContent: "",
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var file *os.File
			var err error
			if tt.fileContent != "" {

				file, err = os.CreateTemp("", "test_config*.json")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(file.Name())

				_, err = file.WriteString(tt.fileContent)
				if err != nil {
					t.Fatalf("Failed to write config data: %v", err)
				}
				file.Close()
			}

			var cfg *config.Config
			if tt.fileContent != "" {
				cfg, err = config.LoadConfig(file.Name())
			} else {

				cfg, err = config.LoadConfig("nonexistent_file.json")
			}

			if (err != nil) != tt.expectedErr {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(cfg, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, cfg)
			}
		})
	}
}
func TestLoadTestConfig(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    *config.ConfigTest
		expectedErr bool
	}{
		{
			name: "Valid Config Test",
			fileContent: `{
				"db": {
					"user": "test user",
					"password": "test password",
					"host": "localhost",
					"name": "test_db"
				}
			}`,
			expected: &config.ConfigTest{
				DB: config.DBConfig{
					User:     "test user",
					Password: "test password",
					Host:     "localhost",
					Name:     "test_db",
				},
			},
			expectedErr: false,
		},
		{
			name: "Invalid JSON",
			fileContent: `{
				"db": {
					"user": "test user",
					"password": "test password",
					"host": "localhost",
					"name": 1234
				}}`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Missing File",
			fileContent: "",
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var file *os.File
			var err error
			if tt.fileContent != "" {

				file, err = os.CreateTemp("", "test_config*.json")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(file.Name())

				_, err = file.WriteString(tt.fileContent)
				if err != nil {
					t.Fatalf("Failed to write config data: %v", err)
				}
				file.Close()
			}

			var cfg *config.ConfigTest
			if tt.fileContent != "" {
				cfg, err = config.LoadTestConfig(file.Name())
			} else {

				cfg, err = config.LoadTestConfig("nonexistent_file.json")
			}

			if (err != nil) != tt.expectedErr {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(cfg, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, cfg)
			}
		})
	}
}
