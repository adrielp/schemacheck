package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

// set default constants for usage messages and default file names
const (
	defaultSchema = "test_data/schema.json"
	schemaUsage   = "A valid JSON schema file to use for validation. Default: schema.json"

	defaultFileName = "test_data/values.json"
	fileUsage       = "A Yaml or JSON file to check against a given schema. Default: values.json (can acceptable multiples)"
)

// Gloval variables for flags and logger
var (
	// Core flag variables
	File   []string
	Schema string

	// Info and Error loggers
	logger    = log.New(os.Stderr, "INFO: ", log.Lshortfile)
	errLogger = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
)

// initialize the flags from the command line and their shorthand counterparts
func init() {
	defaultFile := []string{defaultFileName}
	flag.StringVarP(&Schema, "schema", "s", defaultSchema, schemaUsage)
	flag.StringSliceVarP(&File, "file", "f", defaultFile, fileUsage)
}

// Checks whether a given file is of the supported extension type and if not
// returns false with an error.
// Valid file extensions are currently .yaml, .yml, and .json
func CheckFileIsSupported(file string, fileExt string) (bool, error) {
	// default to false
	fileValid := false

	// supported file extensions to check
	supportedTypes := []string{"yaml", "yml", "json"}

	for _, ext := range supportedTypes {
		if strings.HasSuffix(file, ext) {
			logger.Printf("File: \"%s\" has valid file extension: \"%s\"", file, ext)
			fileValid = true
		}
	}

	if !fileValid {
		return fileValid, errors.New("file type not supported")
	}

	return fileValid, nil

}

func GetFileExt(file string) (string, error) {
	_, fileExt, found := strings.Cut(file, ".")
	if !found {
		return "", errors.New("file separator not found")
	}

	return fileExt, nil
}

func Validate(file string, fileExt string, loadedSchema gojsonschema.JSONLoader) error {
	data, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		errLogger.Panicf("Could not read file: '%s' cleanly.", file)
	}

	if fileExt == "yaml" || fileExt == "yml" {
		data, err = yaml.YAMLToJSON(data)
		if err != nil {
			logger.Panicf("Failed to convert yaml to json in yaml file %s", file)
		}
	}

	documentLoader := gojsonschema.NewBytesLoader(data)

	// Validate the JSON data against the loaded JSON Schema
	result, err := gojsonschema.Validate(loadedSchema, documentLoader)
	if err != nil {
		errLogger.Printf("There was a problem validating %s", file)
		logger.Panicf(err.Error())
	}

	// Check the validity of the result and throw a message is the document is valid or if it's not with errors.
	if result.Valid() {
		logger.Printf("%s is a valid document.\n", file)
	} else {
		logger.Printf("%s is not a valid document...\n", file)
		for _, desc := range result.Errors() {
			errLogger.Printf("--- %s\n", desc)
		}
		return errors.New("document not valid")
	}

	return nil
}

func main() {

	// parse the flags set in the init() function
	flag.Parse()

	// Load schema file before running through and validating the other files to
	// reduce how many times it's loaded.
	schema, err := os.ReadFile(filepath.Clean(Schema))
	if err != nil {
		errLogger.Panicf("Could not read schema file: '%s' cleanly.", Schema)
	}
	loadedSchema := gojsonschema.NewBytesLoader(schema)

	// Iterate through the files declared in the arguments and run validations
	for _, file := range File {
		// Create a specific logger with an ERROR message for easy readability.

		// Print out the values passed on the command line
		logger.Printf("Validating %s file against %s schema...", file, Schema)

		// Get the file extension and error if it failed
		fileExt, err := GetFileExt(file)
		if err != nil {
			errLogger.Panicf(err.Error())
		}

		// Pass the file name and extension to ensure it's a supported file type
		if _, err := CheckFileIsSupported(file, fileExt); err != nil {
			errLogger.Panicf(err.Error())
		}

		Validate(file, fileExt, loadedSchema)
	}
}
