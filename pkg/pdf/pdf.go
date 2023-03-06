package pdf

import (
	"os"
	"path/filepath"
	"time"
	"trikliq-airport-finder/pkg/crypto"
	"trikliq-airport-finder/pkg/logger"

	"go.uber.org/zap"
)

// PdfToTxt is used when we have to process pdf files, to first convert it to txt and then process later, to make it faster
func PdfToTxt(file []byte) (txtFile string, err error) {
	var txtFileOpened []byte
	uuid, _ := crypto.UUID()
	tmpFile := filepath.Join("/tmp", uuid)
	tmpFileTxt := filepath.Join("/tmp", uuid+".txt")

	err = os.WriteFile(tmpFile, file, 0644)
	if err != nil {
		logger.Log.Error("Failed write a file",
			zap.Error(err),
		)
		return
	}
	defer func() {
		os.Remove(tmpFile)
	}()

	args := []string{
		uuid,
		uuid + ".txt",
	}

	_, stderr, err := RunCommand("/tmp", "pdftotext", 10*time.Second, nil, args...)
	if err != nil {
		logger.Log.Error("Failed to bind input data",
			zap.Error(err),
			zap.String("stderr", stderr),
		)
		return

	}
	defer func() {
		os.Remove(tmpFileTxt)
	}()

	txtFileOpened, err = os.ReadFile(tmpFileTxt)
	if err != nil {
		logger.Log.Error("Failed read a file",
			zap.Error(err),
		)
		return
	}

	txtFile = string(txtFileOpened)

	return
}
