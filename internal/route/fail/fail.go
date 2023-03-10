package fail

import (
	"encoding/json"
	"fmt"
	"trikliq-airport-finder/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Error struct {
	Message string `json:"message"`
	Label   string `json:"label"`
}

func SystemError(requestID any) string {
	return fmt.Sprintf("System error happened. Request ID: %s", requestID)
}

func ReturnError(ctx *gin.Context, response model.Response, errMessage string, log *zap.Logger) {
	var (
		raw []byte
		err error
	)

	if errMessage != "" {
		response.Errors = append(response.Errors, errMessage)
	}
	response.Status = false

	raw, err = json.Marshal(response)
	if err != nil {
		log.Error("Failed to marshal data",
			zap.Error(err),
		)

		ctx.JSON(400, response)
		return
	}

	ctx.Data(400, "application/json", raw)
}
