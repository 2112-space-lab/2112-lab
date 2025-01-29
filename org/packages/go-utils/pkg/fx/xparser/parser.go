package xparser

import "encoding/json"

// MapperFunc defines a function type for mapping one type to another
type MapperFunc[T any] func(T) (T, error)

// SerializeStruct converts a struct into a map[string]interface{}, allowing optional mapping
func SerializeStruct[T any](input T, mapper MapperFunc[T]) (map[string]interface{}, error) {
	// If a mapper is provided, transform the input
	if mapper != nil {
		var err error
		input, err = mapper(input)
		if err != nil {
			return nil, err
		}
	}

	// Convert to JSON and then to map[string]interface{}
	data := make(map[string]interface{})
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeserializeStruct converts a map[string]interface{} back into a struct, allowing optional mapping
func DeserializeStruct[T any](data map[string]interface{}, mapper MapperFunc[T]) (T, error) {
	var output T

	// Convert map to JSON and then to struct
	jsonData, err := json.Marshal(data)
	if err != nil {
		return output, err
	}
	err = json.Unmarshal(jsonData, &output)
	if err != nil {
		return output, err
	}

	// If a mapper is provided, apply it
	if mapper != nil {
		return mapper(output)
	}

	return output, nil
}
