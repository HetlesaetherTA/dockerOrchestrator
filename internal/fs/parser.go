package fs

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

func ConstructFromYAML(object any, path string) error {
	rv := reflect.ValueOf(object)

	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("object must be a non-nil pointer")
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("Could not read file from %s: %w", path, err)
	}

	err = yaml.Unmarshal([]byte(data), object)

	return err
}
