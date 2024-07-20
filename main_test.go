package main

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestReadEnvFile(t *testing.T) {
	// Create a temporary .env file
	content := []byte("KEY1=value1\nKEY2=value2\n# Comment\nKEY3=value3 with spaces")
	tmpfile, err := os.CreateTemp("", "test.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test reading the .env file
	secrets, err := readEnvFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("readEnvFile() error = %v", err)
	}

	expectedKeys := []string{"KEY1", "KEY2", "KEY3"}
	for _, key := range expectedKeys {
		if _, ok := secrets[key]; !ok {
			t.Errorf("Expected key %s not found in secrets", key)
		}
	}

	if len(secrets) != 3 {
		t.Errorf("Expected 3 secrets, got %d", len(secrets))
	}
}

func TestCreateSecretYAML(t *testing.T) {
	secrets := map[string]string{
		"KEY1": "dmFsdWUx",
		"KEY2": "dmFsdWUy",
	}

	// Test without namespace
	yamlContent, err := createSecretYAML(secrets, "name", "")
	if err != nil {
		t.Fatalf("createSecretYAML() error = %v", err)
	}

	var secret Secret
	err = yaml.Unmarshal(yamlContent, &secret)
	if err != nil {
		t.Fatalf("Error unmarshaling YAML: %v", err)
	}

	if secret.ApiVersion != "v1" {
		t.Errorf("Expected ApiVersion v1, got %s", secret.ApiVersion)
	}
	if secret.Kind != "Secret" {
		t.Errorf("Expected Kind Secret, got %s", secret.Kind)
	}
	if secret.Metadata["name"] != "name" {
		t.Errorf("Expected name env-secrets, got %s", secret.Metadata["name"])
	}
	if _, ok := secret.Metadata["namespace"]; ok {
		t.Error("Namespace should not be present")
	}

	// Test with namespace
	yamlContent, err = createSecretYAML(secrets, "name", "my-namespace")
	if err != nil {
		t.Fatalf("createSecretYAML() error = %v", err)
	}

	err = yaml.Unmarshal(yamlContent, &secret)
	if err != nil {
		t.Fatalf("Error unmarshaling YAML: %v", err)
	}

	if secret.Metadata["namespace"] != "my-namespace" {
		t.Errorf("Expected namespace my-namespace, got %s", secret.Metadata["namespace"])
	}
}
