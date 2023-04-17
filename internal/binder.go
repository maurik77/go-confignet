package internal

import (
	"log"
	"reflect"
	"strconv"
	"time"
)

func Unmarshal(target interface{}, value string, parts ...string) {
	reflectedType := reflect.ValueOf(target).Elem()
	fillObject(reflectedType, value, parts...)
}

func fillObject(parent reflect.Value, value string, parts ...string) {
	log.Printf("fillObject -> parts: %v, value: %v", parts, value)

	if parent.Kind() != reflect.Struct {
		log.Printf("Configuration:Parent is not a valid struct %v. Parts %v", parent, parts)
		return
	}

	fieldName := parts[0]
	nestedField := parent.FieldByName(fieldName)

	if !nestedField.IsValid() {
		log.Printf("Configuration:Unable to find field %v in the object %v", fieldName, nestedField)
		return
	}

	nestedField = checkValueOfPointer(nestedField)

	switch {
	case len(parts) == 1: // property
		fillField(nestedField, value, 0)
	case nestedField.Kind() == reflect.Slice:
		fillSlice(nestedField, value, parts...)
	case nestedField.Kind() == reflect.Array:
		fillArray(nestedField, value, parts...)
	case nestedField.Kind() == reflect.Map:
		fillMap(nestedField, value, parts...)
	default: // nested object
		fillObject(nestedField, value, parts[1:]...)
	}
}

func checkValueOfPointer(field reflect.Value) reflect.Value {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			defaultValue := reflect.New(field.Type().Elem()).Elem()
			field.Set(defaultValue.Addr())
		}
		return field.Elem()
	}

	return field
}

func fillMap(nestedField reflect.Value, value string, parts ...string) {
	log.Printf("fillMap -> parts: %v, value: %v", parts, value)

	mapKeyType := nestedField.Type().Key()
	if mapKeyType.Kind() == reflect.Struct {
		log.Printf("this version is not able to manage maps having struct as key")
		return
	}

	mapValueType := nestedField.Type().Elem()
	key := parts[1]
	log.Printf("  	mapKeyType: %v, mapValueType: %v, key: %v, value path: %v, value: %v", mapKeyType, mapValueType, key, parts[2:], value)

	if nestedField.IsNil() {
		elemMap := reflect.MakeMap(reflect.MapOf(mapKeyType, mapValueType))
		nestedField.Set(elemMap)
	}

	keyField := reflect.New(mapKeyType)
	fillField(keyField, key, -1)
	valueField := reflect.New(nestedField.Type().Elem()).Elem()
	existingValue := nestedField.MapIndex(keyField.Elem())

	if existingValue.IsValid() {
		valueField.Set(existingValue)
	}

	if len(parts) == 2 {
		fillField(valueField, value, -1)
	} else {
		valueFieldElem := checkValueOfPointer(valueField)
		fillObject(valueFieldElem, value, parts[2:]...)
	}

	nestedField.SetMapIndex(keyField.Elem(), valueField)
}

func fillArray(nestedField reflect.Value, value string, parts ...string) {
	log.Printf("fillArray -> parts: %v, value: %v", parts, value)

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Configuration:Unable to parse index %v for path %v in the object %v", parts[1], parts, nestedField)
		return
	}

	if index >= nestedField.Len() {
		log.Printf("Configuration:Unable to assign value %v to the field %v with index %v because is out of range. (Array length  %v)", value, parts[0], index, nestedField.Len())
		return
	}

	if len(parts) == 2 {
		fillField(nestedField, value, index)
	} else {
		valueFieldElem := checkValueOfPointer(nestedField.Index(index))
		fillObject(valueFieldElem, value, parts[2:]...)
	}
}

func fillSlice(nestedField reflect.Value, value string, parts ...string) {
	log.Printf("fillSlice -> parts: %v, value: %v", parts, value)

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Configuration:Unable to parse index %v for path %v in the object %v", parts[1], parts, nestedField)
		return
	}

	if nestedField.IsNil() {
		elemType := nestedField.Type().Elem()
		elemSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
		nestedField.Set(elemSlice)
	}

	if index >= nestedField.Len() {
		newElements := make([]reflect.Value, index+1-nestedField.Len())
		for i := range newElements {
			newElements[i] = reflect.New(nestedField.Type().Elem()).Elem()
		}
		nestedField.Set(reflect.Append(nestedField, newElements...))
	}

	if len(parts) == 2 {
		fillField(nestedField, value, index)
	} else {
		valueFieldElem := checkValueOfPointer(nestedField.Index(index))
		fillObject(valueFieldElem, value, parts[2:]...)
	}
}

func fillField(field reflect.Value, value string, index int) {
	log.Printf("fillField -> index: %v, value: %v, kind: %v", index, value, field.Kind())

	switch field.Kind() {
	case reflect.Array:
		item := field.Index(index)
		fillField(item, value, -1)
		return
	case reflect.Slice:
		item := field.Index(index)
		fillField(item, value, -1)
		return
	case reflect.Ptr:
		if field.IsNil() {
			defaultValue := reflect.New(field.Type().Elem()).Elem()
			field.Set(defaultValue.Addr())
		}
		fillField(field.Elem(), value, -1)
	default:
		switch field.Interface().(type) {
		case string:
			field.SetString(value)
		case int, int8, int16, int32, int64:
			if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
				if field.CanSet() {
					field.SetInt(intValue)
				}
			} else {
				log.Printf("Configuration:Unable to parse Int the value %v", value)
			}
		case uint, uint8, uint16, uint32, uint64:
			if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
				if field.CanSet() {
					field.SetUint(uintValue)
				}
			} else {
				log.Printf("Configuration:Unable to parse Uint the value %v", value)
			}
		case bool:
			if boolValue, err := strconv.ParseBool(value); err == nil {
				if field.CanSet() {
					field.SetBool(boolValue)
				}
			} else {
				log.Printf("Configuration:Unable to parse Bool the value %v", value)
			}
		case time.Time:
			if timeValue, err := time.Parse(time.RFC3339Nano, value); err == nil {
				if field.CanSet() {
					field.Set(reflect.ValueOf(timeValue))
				}
			} else {
				log.Printf("Configuration:Unable to parse Time (format: %v) the value %v", time.RFC3339Nano, value)
			}
		}
	}
}
