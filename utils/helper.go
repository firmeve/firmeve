package utils

func InSlice(value interface{}, slices []interface{}) bool {
	for _, v := range slices {
		if v == value {
			return true
		}
	}

	return false
}

func InMap(key interface{}, items map[interface{}]interface{}) bool {
	if _, ok := items[key]; ok {
		return true
	}

	return false
}
