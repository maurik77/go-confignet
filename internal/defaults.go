package internal

import "reflect"

// ApplyDefaults walks the struct pointed to by target and sets any field
// tagged with `default:"<value>"` to that value when the field currently
// holds its zero value. Provider-supplied values are applied afterwards by
// the caller, so they always take precedence over tag defaults.
func ApplyDefaults(target interface{}) {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	applyDefaultsToStruct(rv.Elem())
}

func applyDefaultsToStruct(v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Apply the default only for scalar/pointer fields (not structs themselves).
		if defaultVal, ok := fieldType.Tag.Lookup("default"); ok &&
			field.Kind() != reflect.Struct &&
			field.IsZero() {
			fillField(field, defaultVal, -1)
		}

		// Recurse into nested structs and already-allocated pointer-to-structs.
		switch field.Kind() {
		case reflect.Struct:
			applyDefaultsToStruct(field)
		case reflect.Ptr:
			if !field.IsNil() {
				if elem := field.Elem(); elem.Kind() == reflect.Struct {
					applyDefaultsToStruct(elem)
				}
			}
		}
	}
}
