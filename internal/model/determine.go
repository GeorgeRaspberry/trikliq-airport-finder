package model

func Determine(entity string) interface{} {
	switch entity {
	case "user":
		return &User{}
	}

	return nil
}

func DetermineArray(entity string) interface{} {
	switch entity {
	case "user":
		return &[]User{}
	}

	return nil
}
