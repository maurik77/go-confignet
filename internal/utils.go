package internal

import "fmt"

//LoadProperties loads value map in a slice map[string]string
func LoadProperties(separator string, valueMap map[string]interface{}) map[string]string {
	target := map[string]string{}
	loadProperties("", separator, valueMap, target)
	return target
}

func loadProperties(parent string, separator string, valueMap map[string]interface{}, data map[string]string) {
	for key, value := range valueMap {
		if parent != "" {
			key = fmt.Sprintf("%v%v%v", parent, separator, key)
		}

		switch v := value.(type) {
		default:
			data[key] = fmt.Sprint(value)
		case []interface{}:
			loadArray(key, separator, v, data)
		case map[string]interface{}:
			loadProperties(key, separator, v, data)
		}
	}
}

func loadArray(key string, separator string, slice []interface{}, data map[string]string) {
	for index, value := range slice {
		indexKey := fmt.Sprintf("%v%v%v", key, separator, index)
		switch v := value.(type) {
		default:
			data[indexKey] = fmt.Sprint(value)
		case []interface{}:
			loadArray(indexKey, separator, v, data)
		case map[string]interface{}:
			loadProperties(indexKey, separator, v, data)
		}
	}
}
