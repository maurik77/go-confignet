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
	// log.Printf("fillObject -> parts: %v, value: %v", parts, value)

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

	if nestedField.Kind() == reflect.Ptr {
		if nestedField.IsNil() {
			defaultValue := reflect.New(nestedField.Type().Elem()).Elem()
			nestedField.Set(defaultValue.Addr())
		}
		nestedField = nestedField.Elem()
	}

	switch {
	case len(parts) == 1: // property
		fillField(nestedField, value, 0)
	case nestedField.Kind() == reflect.Slice:
		fillSlice(fieldName, nestedField, value, parts...)
	case nestedField.Kind() == reflect.Array:
		fillArray(fieldName, nestedField, value, parts...)
	// case nestedField.Kind() == reflect.Map:
	// 	log.Printf("fillObject::Map: parts %v value %v", parts, value)
	default: // nested object
		fillObject(nestedField, value, parts[1:]...)
	}
}

func fillArray(fieldName string, nestedField reflect.Value, value string, parts ...string) {
	// log.Printf("fillArray -> fieldName: %v, parts: %v, value: %v", fieldName, parts, value)

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Configuration:Unable to parse index %v for field %v in the object %v", parts[1], fieldName, nestedField)
		return
	}

	if index >= nestedField.Len() {
		log.Printf("Configuration:Unable to assign value %v to the field %v with index %v because is out of range. (Array length  %v)", value, fieldName, index, nestedField.Len())
		return
	}

	if len(parts) == 2 {
		fillField(nestedField, value, index)
	} else {
		fillObject(nestedField.Index(index), value, parts[2:]...)
	}
}

func fillSlice(fieldName string, nestedField reflect.Value, value string, parts ...string) {
	// log.Printf("fillSlice -> fieldName: %v, parts: %v, value: %v", fieldName, parts, value)

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Configuration:Unable to parse index %v for field %v in the object %v", parts[1], fieldName, nestedField)
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
		fillObject(nestedField.Index(index), value, parts[2:]...)
	}
}

func fillField(field reflect.Value, value string, index int) {
	// log.Printf("fillField -> index: %v, value: %v", index, value)

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
				field.SetInt(intValue)
			} else {
				log.Printf("Configuration:Unable to parse Int the value %v", value)
			}
		case uint, uint8, uint16, uint32, uint64:
			if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
				field.SetUint(uintValue)
			} else {
				log.Printf("Configuration:Unable to parse Uint the value %v", value)
			}
		case bool:
			if boolValue, err := strconv.ParseBool(value); err == nil {
				field.SetBool(boolValue)
			} else {
				log.Printf("Configuration:Unable to parse Bool the value %v", value)
			}
		case time.Time:
			if timeValue, err := time.Parse(time.RFC3339Nano, value); err == nil {
				field.Set(reflect.ValueOf(timeValue))
			} else {
				log.Printf("Configuration:Unable to parse Time (format: %v) the value %v", time.RFC3339Nano, value)
			}
		}
	}
}
