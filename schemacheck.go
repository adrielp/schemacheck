package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

// Set default constants for flag usage messages.
const (
	schemaUsage              = "A valid JSON schema file to use for validation. Default: schema.json"
	fileUsage                = "A Yaml or JSON file to check against a given schema. Default: values.json (can acceptable multiples)"
	versionUsage             = "Prints out the version of schemacheck"
	ignoreValidationErrUsage = "Ignores when a document is not valid but provides a warning."
)

// Core variables for flag pointers and info, warning, and error loggers.
var (
	// Core flag variables
	File                []string
	Schema              string
	IgnoreValidationErr bool
	VersionFlag         bool

	// version is set through ldflags by GoReleaser upon build, taking in the most recent tag
	// and appending -snapshot in the event that --snapshot is set in GoReleaser.
	version string

	// Info, warning, and error loggers.
	logger     = log.New(os.Stderr, "INFO: ", log.Lshortfile)
	warnLogger = log.New(os.Stderr, "WARN: ", log.Lshortfile)
	errLogger  = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
)

// Initialize the flags from the command line and their shorthand counterparts.
func init() {
	flag.StringVarP(&Schema, "schema", "s", "", schemaUsage)
	flag.StringSliceVarP(&File, "file", "f", []string{}, fileUsage)
	flag.BoolVar(&IgnoreValidationErr, "ignore-val-err", false, ignoreValidationErrUsage)
	flag.BoolVarP(&VersionFlag, "version", "v", false, versionUsage)
}

// Check whether or not a required flag like file and schema is set and return true or false.
func CheckForEmptyArg() bool {
	schemaArgEmpty := true
	fileArgEmpty := true
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "schema" {
			if f.Changed {
				schemaArgEmpty = false
			}
		} else if f.Name == "file" {
			if f.Changed {
				fileArgEmpty = false
			}
		}
	})
	if schemaArgEmpty || fileArgEmpty {
		return true
	}
	return false
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
		errLogger.Fatalf("Could not read file: '%s' cleanly.", file)
	}

	if fileExt == "yaml" || fileExt == "yml" {
		data, err = yaml.YAMLToJSON(data)
		if err != nil {
			logger.Fatalf("Failed to convert yaml to json in yaml file %s", file)
		}
	}

	documentLoader := gojsonschema.NewBytesLoader(data)

	// Validate the JSON data against the loaded JSON Schema
	result, err := gojsonschema.Validate(loadedSchema, documentLoader)
	if err != nil {
		errLogger.Printf("There was a problem validating %s", file)
		logger.Fatalf(err.Error())
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

	// If version flag is set, output version of app and exit
	if VersionFlag {
		fmt.Printf("schemacheck version: %s\n", version)
		os.Exit(0)
	}

	// Check to ensure required flags aren't empty
	missingArgs := CheckForEmptyArg()
	if missingArgs {
		fmt.Fprintf(os.Stderr, "Usage of schemacheck\n")
		flag.PrintDefaults()
		errLogger.Fatal("One or more missing args not set.")
	}

	// Load schema file before running through and validating the other files to
	// reduce how many times it's loaded.
	schema, err := os.ReadFile(filepath.Clean(Schema))
	if err != nil {
		errLogger.Fatalf("Could not read schema file: '%s' cleanly.", Schema)
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
			errLogger.Fatalf(err.Error())
		}

		// Pass the file name and extension to ensure it's a supported file type
		if _, err := CheckFileIsSupported(file, fileExt); err != nil {
			errLogger.Fatal(err.Error())
		}

		// Validate against the schema and if IgnoreValidationErr is set, exit with a warning.
		if err := Validate(file, fileExt, loadedSchema); err != nil {
			if IgnoreValidationErr {
				warnLogger.Printf("Ignoring validation error.")
				os.Exit(0)
			}
			errLogger.Fatal(err.Error())
		}
	}
}
