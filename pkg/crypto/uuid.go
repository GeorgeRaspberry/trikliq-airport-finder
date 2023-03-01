package crypto

import "github.com/google/uuid"

// UUID will return new UUID
func UUID() (string, error) {
	uuidValue, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuidValue.String(), err
}
