package crypto

import (
	"crypto/sha256"
	"fmt"
)

func SHA256(input interface{}) (output string) {
	value := fmt.Sprintf("%+v", input)

	h := sha256.New()
	h.Write([]byte(value))

	bs := string(h.Sum(nil))
	output = fmt.Sprintf("%x", bs)

	return
}

func HashPassword(password string) string {
	salt, _ := UUID()
	password = fmt.Sprintf("%s.%s", salt, password)
	password = fmt.Sprintf("%s.%s", salt, SHA256(password))
	return password
}
