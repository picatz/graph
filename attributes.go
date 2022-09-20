package graph

import "fmt"

// Attributes are named values that can be associated with a node or subgraph.
type Attributes map[string]any

// UseAttribute is a helper function to use a named attribute of a specific type.
func UseAttribute[T any](attrs Attributes, name string, fn func(T)) error {
	v, ok := attrs[name]
	if !ok {
		return fmt.Errorf("graph node attribute %q doesn't exist", name)
	}
	vt, ok := v.(T)
	if !ok {
		return fmt.Errorf("graph node attribute %q is of type %T not %T", name, v, vt)
	}
	fn(vt)
	return nil
}
