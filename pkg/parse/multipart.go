package parse

import (
	"io"
	"mime/multipart"
	"trikliq-airport-finder/internal/model"
)

func ReadMultipartFiles(files []*multipart.FileHeader) ([]model.MultipartFile, error) {
	rawFiles := make([]model.MultipartFile, 0)

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return nil, err
		}

		raw, err := io.ReadAll(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		f.Close()

		mp := model.MultipartFile{
			Filename: file.Filename,
			Size:     file.Size,
			Header:   file.Header,
			Content:  raw,
		}

		rawFiles = append(rawFiles, mp)
	}

	return rawFiles, nil
}
