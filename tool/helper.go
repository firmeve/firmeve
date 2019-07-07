package tool

func InSlice(value interface{}, slices []interface{}) bool {
	for _, v := range slices {
		if v == value {
			return true
		}
	}

	return false
}
