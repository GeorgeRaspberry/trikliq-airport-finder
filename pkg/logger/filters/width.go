package filters

import "bytes"

// MinWidth will make string with minimal width
func MinWidth(in, separator string, min int) string {
	diff := min - len(in)

	if diff > 0 {
		var buffer bytes.Buffer
		buffer.WriteString(in)

		for i := 0; i < diff; i++ {
			buffer.WriteString(separator)
		}
		in = buffer.String()
	}

	return in
}
