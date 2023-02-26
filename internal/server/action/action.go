package action

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// encodeResponse in proper formatted output
func encodeResponse(res interface{}) (out string, err error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)

	// NOTES: escape XSS
	// we can do this manually by escaping: < > ' " &
	enc.SetEscapeHTML(true)

	enc.SetIndent("", "\t")
	err = enc.Encode(res)
	out = b.String()

	return
}

//go:generate easytags $GOFILE json:camel,form:camel
type ResponseError struct {
	Error string `json:"error" form:"error"`
}

// Response performs response on stdout
func Response(ctx *gin.Context, res interface{}, bsErrors []error) {
	ctx.Header("Content-Type", "application/json; charset=utf-8")

	if len(bsErrors) > 0 {
		responseErrors := make([]ResponseError, 0)

		for _, err := range bsErrors {
			if err == nil {
				continue
			}

			responseError := ResponseError{
				Error: err.Error(),
			}
			responseErrors = append(responseErrors, responseError)
		}

		obj, _ := encodeResponse(responseErrors)
		ctx.String(401, obj)
		return
	}

	obj, _ := encodeResponse(res)
	statusCode := ctx.Writer.Status()
	ctx.Data(statusCode, "application/json; charset=utf-8", []byte(obj))
}
