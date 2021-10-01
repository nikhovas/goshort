package utils

func UnwrapFieldOrDefault(dict map[string]interface{}, key string, defaultValue interface{}) interface{} {
	value, ok := dict[key]
	if !ok {
		return defaultValue
	} else {
		return value
	}
}
