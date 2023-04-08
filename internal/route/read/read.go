package read

import (
	"trikliq-airport-finder/internal/model"
	"trikliq-airport-finder/internal/route/fail"
	"trikliq-airport-finder/internal/server/router"
	"trikliq-airport-finder/pkg/logger"
	"trikliq-airport-finder/pkg/parse"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ReadHandler(ctx *gin.Context) {

	var (
		//request  = model.Request{}
		response = model.Response{}
	)

	requestID, _ := ctx.Get("id")

	log := logger.Log.WithOptions(zap.Fields(
		zap.Any("requestID", requestID),
	))

	log.Info("read started")

	switch ctx.ContentType() {
	case "multipart/form-data":

		form, err := ctx.MultipartForm()
		if err != nil {
			log.Error("form not found",
				zap.Error(err),
			)
			fail.ReturnError(ctx, response, err.Error(), log)
			return
		}

		// if no documents are here, skip the upload part
		if len(form.File) == 0 {
			log.Error("no files found in form",
				zap.Any("form", *form),
			)

			fail.ReturnError(ctx, response, "no files found in form", log)
			return
		}

		files := make([]model.MultipartFile, 0)
		for _, fileHeaders := range form.File {
			rawFiles, err := parse.ReadMultipartFiles(fileHeaders)
			if err != nil {
				log.Error("error while reading multipart file(s)",
					zap.Error(err),
				)
				continue
			}
			files = append(files, rawFiles...)
		}

		result := make(map[string]any, 0)

		for _, file := range files {
			result[file.Filename] = parse.Parse(file.Content, log)
		}

		response.Data = result

	default:

		fail.ReturnError(ctx, response, "content-type is not multipart/form-data", log)
		return
	}

	log.Info("create finished")

	response.Status = true
	ctx.JSON(200, response)
}

func init() {
	router.Router.Handle("POST", "/read", ReadHandler)
}
