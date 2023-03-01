package security

import (
	"bytes"
	"fmt"
	"strings"

	libTransform "trikliq-airport-finder/pkg/transform"
)

//lint:ignore GLOBAL this is okay
var (
	monitoringKeys = []string{
		"apikey",
		"api_",
		"private",
		"secret",
		"secure",
		"password",
		"session",
		"token",
		"awsaccessid",
		"awssecretkey",
		"auth",
		"aws",
		"credential",
		"jwt",
		"authorization",
		"accesskey",
	}
)

// ObfuscateField takes key and value, if key is in our list than value will be obfuscated and key wil be marked with red color
func ObfuscateField(key string, val interface{}) (outKey string, outVal interface{}) {
	outKey = key

	keySanitized := strings.ToLower(key)
	isCandidate := false
	for _, monitoringKey := range monitoringKeys {
		monitoringKeySanitized := strings.ToLower(monitoringKey)

		if strings.Contains(keySanitized, monitoringKeySanitized) {
			isCandidate = true
			break
		}
	}

	if isCandidate {
		outKey = fmt.Sprintf("%s*", key)
		libTransform.AnyToSliceOfString(val)
		val = fmt.Sprintf("%v", val)

		var (
			buffer bytes.Buffer
			total  int
		)
		total = len(val.(string))

		for pos, char := range val.(string) {
			if pos < 4 {
				chr := fmt.Sprintf("%c", char)
				buffer.WriteString(chr)
			} else {
				if total > 20 {
					if pos > total-4 {
						chr := fmt.Sprintf("%c", char)
						buffer.WriteString(chr)
					} else {
						buffer.WriteString("*")
					}
				} else {
					buffer.WriteString("*")
				}
			}
		}
		outVal = buffer.String()
	} else {
		outVal = fmt.Sprintf("%v", val)
		outVal, _ = outVal.(string)
	}

	return
}
