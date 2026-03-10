package internal

import (
	"reflect"
	"sync"
)

const TagName = "confignet"

// fieldCache maps reflect.Type → (key string → field index).
// Built once per struct type and reused on every subsequent lookup.
var fieldCache sync.Map // map[reflect.Type]map[string]int

// fieldByKey returns the struct field of v whose confignet tag (or, if absent,
// whose field name) matches key. The second return value is false when no
// matching field exists.
func fieldByKey(v reflect.Value, key string) (reflect.Value, bool) {
	index, ok := cachedFieldIndex(v.Type(), key)
	if !ok {
		return reflect.Value{}, false
	}
	return v.Field(index), true
}

// typeFieldIndex returns the field index inside type t whose confignet tag
// (or field name) matches key, using the per-type cache.
func typeFieldIndex(t reflect.Type, key string) (int, bool) {
	return cachedFieldIndex(t, key)
}

func cachedFieldIndex(t reflect.Type, key string) (int, bool) {
	if v, ok := fieldCache.Load(t); ok {
		m := v.(map[string]int)
		idx, found := m[key]
		return idx, found
	}

	m := buildIndex(t)
	fieldCache.Store(t, m)
	idx, found := m[key]
	return idx, found
}

// buildIndex scans all exported fields of t and maps each field's lookup key
// (confignet tag value, falling back to field name) to its index.
func buildIndex(t reflect.Type) map[string]int {
	m := make(map[string]int, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := f.Tag.Get(TagName)
		if key == "" {
			key = f.Name
		}
		m[key] = i
	}
	return m
}
