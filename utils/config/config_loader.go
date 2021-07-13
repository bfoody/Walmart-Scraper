package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// A ConfigLoadError is thrown when the environment is missing one or more variables.
type ConfigLoadError struct {
	MissingFields []string
}

// Error displays the ConfigLoadError's missing variables.
func (c *ConfigLoadError) Error() string {
	return fmt.Sprintf("environment missing variables: %s", strings.Join(c.MissingFields, ", "))
}

// LoadConfigFromEnv takes a struct ptr as input and uses reflection to load fields in
// from environment variables, using the `env` struct tag to determine lookup names.
//
// Returns an error if any variables are not found.
func LoadConfigFromEnv(configStruct interface{}) error {
	elem := reflect.TypeOf(configStruct).Elem()
	if elem.Kind() != reflect.Struct {
		panic("LoadConfigFromEnv must be passed a pointer to a struct")
	}

	val := reflect.ValueOf(configStruct).Elem()

	missingFields := []string{}
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		name, ok := field.Tag.Lookup("env")
		if !ok {
			name = field.Name
		}

		envVal := os.Getenv(name)
		if envVal == "" {
			missingFields = append(missingFields, name)
			continue
		}

		val.Field(i).SetString(envVal)
	}

	if len(missingFields) > 0 {
		return &ConfigLoadError{
			missingFields,
		}
	}

	return nil
}
