package transform

// InSlice checks if there is an element in slice
func InSlice(text string, slice []string) bool {
	for _, i := range slice {
		if i == text {
			return true
		}
	}

	return false
}

func InSliceOfMap(data []map[string]string, target string) bool {
	for _, sample := range data {
		if sample["iata"] == target {
			return true
		}
	}

	return false
}
