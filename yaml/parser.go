package yaml

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

var (
	ErrLoading   = errors.New("error merging YAML files")
	ErrMarshal   = errors.New("error marshalling YAML")
	ErrMerging   = errors.New("error merging YAML")
	ErrUnmarshal = errors.New("error unmarshalling YAML")
)

// RecursiveMerge merges src into dest recursively.
func RecursiveMerge(dest, src map[string]interface{}) map[string]interface{} {
	for key, srcValue := range src {
		if destValue, exists := dest[key]; exists {
			// If both values are maps, merge them recursively
			destMap, ok1 := destValue.(map[string]interface{})
			srcMap, ok2 := srcValue.(map[string]interface{})

			if ok1 && ok2 {
				dest[key] = RecursiveMerge(destMap, srcMap)

				continue
			}
		}

		// Otherwise, overwrite dest with src
		dest[key] = srcValue
	}

	return dest
}

// LoadYAMLFile loads a YAML file and unmarshals it into a map.
func LoadYAMLFile(filepath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.Join(ErrLoading, err)
	}

	var result map[string]interface{}

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, errors.Join(
			ErrUnmarshal,
			fmt.Errorf("file %s", filepath), //nolint:err113 // Expose file name
			err,
		)
	}

	return result, nil
}

// MergeYAMLFiles merges YAML files into one map.
func MergeYAMLFiles(files ...string) (map[string]interface{}, error) {
	merged := make(map[string]interface{})

	for _, file := range files {
		if file == "" {
			continue // Skip empty file paths
		}

		content, err := LoadYAMLFile(file)
		if err != nil {
			return nil, errors.Join(ErrMerging, err)
		}

		merged = RecursiveMerge(merged, content)
	}

	return merged, nil
}

// Unmarshal unmarshals a list of YAML files into a given struct.
func Unmarshal[T any](files ...string) (*T, error) {
	var t T

	// Merge all YAML files
	m, err := MergeYAMLFiles(files...)
	if err != nil {
		return nil, err
	}

	// Marshal the merged result
	mergedYAML, err := yaml.Marshal(m)
	if err != nil {
		return nil, errors.Join(ErrMarshal, err)
	}

	// Unmarshal back into the target struct
	if err := yaml.Unmarshal(mergedYAML, &t); err != nil {
		return nil, errors.Join(ErrUnmarshal, err)
	}

	return &t, nil
}
