package graph

import "fmt"

// Attributes are named values that can be associated with a node or subgraph.
type Attributes map[string]any

// UseAttribute is a helper function to use a named attribute of a specific type.
func UseAttribute[T any](attrs Attributes, name string, fn func(T)) error {
	v, ok := attrs[name]
	if !ok {
		return fmt.Errorf("graph attribute %q doesn't exist", name)
	}
	vt, ok := v.(T)
	if !ok {
		return fmt.Errorf("graph attribute %q is of type %T not %T", name, v, vt)
	}
	fn(vt)
	return nil
}

// GetAttributes is a helper function to return a named attribute of a specific type.
func GetAttribute[T any](attrs Attributes, name string) (T, error) {
	var (
		v   T
		err error
	)
	err = UseAttribute(attrs, name, func(value T) {
		v = value
	})
	if err != nil {
		return v, fmt.Errorf("failed to get attribute: %w", err)
	}
	return v, nil
}

// SetAttribute is a helper function to set a named attribute.
func SetAttribute[T any](attrs Attributes, name string, value T) {
	attrs[name] = value
}

// DeleteAttribute is a helper function to remove a named attribute.
func DeleteAttribute[T any](attrs Attributes, name string) {
	delete(attrs, name)
}
