package internal

import (
	"reflect"
	"strconv"
)

// IsValidPath reports whether the sequence of parts navigates to a known
// field (or element) inside the type t, respecting confignet struct tags.
// It is used by BindStrict to detect configuration keys that have no
// matching struct field.
func IsValidPath(t reflect.Type, parts []string) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if len(parts) == 0 {
		return true
	}

	part := parts[0]
	rest := parts[1:]

	switch t.Kind() {
	case reflect.Struct:
		idx, ok := typeFieldIndex(t, part)
		if !ok {
			return false
		}
		return IsValidPath(t.Field(idx).Type, rest)
	case reflect.Slice, reflect.Array:
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
		return IsValidPath(t.Elem(), rest)
	case reflect.Map:
		// part is the map key — always valid; validate rest against value type
		return IsValidPath(t.Elem(), rest)
	default:
		// scalar — no further parts expected
		return len(rest) == 0
	}
}
