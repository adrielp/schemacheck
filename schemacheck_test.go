package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

// ===============================================
// Tests against the CheckFileIsSupported function for various file types
func TestCheckFileIsSupportedYaml(t *testing.T) {
	file := "test_data/values.yaml"
	valid, err := CheckFileIsSupported(file, "yaml")
	if err != nil {
		t.Fatalf("An error occured during validation of file that should not have occurred:\n %s", err)
	}
	assert.True(t, valid)
}

func TestCheckFileIsSupportedYml(t *testing.T) {
	file := "test_data/values.yml"
	valid, err := CheckFileIsSupported(file, "yml")
	if err != nil {
		t.Fatalf("An error occured during validation of file that should not have occurred:\n %s", err)
	}
	assert.True(t, valid)
}

func TestCheckFileIsSupportedJSON(t *testing.T) {
	file := "test_data/values.json"
	valid, err := CheckFileIsSupported(file, "json")
	if err != nil {
		t.Fatalf("An error occured during validation of file that should not have occurred:\n %s", err)
	}
	assert.True(t, valid)
}

func TestCheckFileIsSupportedTxt(t *testing.T) {
	file := "test_data/values.txt"
	_, err := CheckFileIsSupported(file, "txt")
	assert.Error(t, err)
}

// ===============================================
// Tests GetFileExt function for various filetypes
func TestGetFileExtYaml(t *testing.T) {
	file := "test_data/values.yaml"
	fileExt, _ := GetFileExt(file)
	assert.Equal(t, "yaml", fileExt)
}

func TestGetFileExtYml(t *testing.T) {
	file := "test_data/values.yml"
	fileExt, _ := GetFileExt(file)
	assert.Equal(t, "yml", fileExt)
}

func TestGetFileExtJSON(t *testing.T) {
	file := "test_data/values.json"
	fileExt, _ := GetFileExt(file)
	assert.Equal(t, "json", fileExt)
}

func TestGetFileExtNoSeparator(t *testing.T) {
	file := "test_data/noseparator"
	_, err := GetFileExt(file)
	assert.Error(t, err)
}

// ===============================================
// Tests Validate against test data files
func TestValidateValidYaml(t *testing.T) {
	file := "test_data/values.yaml"
	fileExt := "yaml"
	schema, err := os.ReadFile(filepath.Clean("test_data/schema.json"))
	if err != nil {
		errLogger.Panicf("Could not read schema file: '%s' cleanly.", Schema)
	}
	loadedSchema := gojsonschema.NewBytesLoader(schema)
	assert.NoError(t, Validate(file, fileExt, loadedSchema))
}

func TestValidateValidJSON(t *testing.T) {
	file := "test_data/values.json"
	fileExt := "yaml"
	schema, err := os.ReadFile(filepath.Clean("test_data/schema.json"))
	if err != nil {
		errLogger.Panicf("Could not read schema file: '%s' cleanly.", Schema)
	}
	loadedSchema := gojsonschema.NewBytesLoader(schema)
	assert.NoError(t, Validate(file, fileExt, loadedSchema))
}

func TestValidateInvalidYaml(t *testing.T) {
	file := "test_data/invalid.yaml"
	fileExt := "yaml"
	schema, err := os.ReadFile(filepath.Clean("test_data/schema.json"))
	if err != nil {
		errLogger.Panicf("Could not read schema file: '%s' cleanly.", Schema)
	}
	loadedSchema := gojsonschema.NewBytesLoader(schema)
	assert.Error(t, Validate(file, fileExt, loadedSchema))
}

func TestValidateInvalidJSON(t *testing.T) {
	file := "test_data/invalid.json"
	fileExt := "yaml"
	schema, err := os.ReadFile(filepath.Clean("test_data/schema.json"))
	if err != nil {
		errLogger.Panicf("Could not read schema file: '%s' cleanly.", Schema)
	}
	loadedSchema := gojsonschema.NewBytesLoader(schema)
	assert.Error(t, Validate(file, fileExt, loadedSchema))
}
