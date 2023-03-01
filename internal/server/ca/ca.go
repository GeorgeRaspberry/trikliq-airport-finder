package ca

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	libIO "trikliq-airport-finder/pkg/io"

	"go.uber.org/zap"
)

//lint:ignore GLOBAL this is okay
var (
	execFile  string
	homeDir   string
	logger    *zap.Logger
	hostnames = []string{
		"localhost",
		"127.0.0.1",
		"::1",
	}
)

func readDir(homeDir string) (res []string, err error) {
	res = make([]string, 0)

	f, err := os.Open(homeDir)
	if err != nil {
		logger.Error("error in listing directory",
			zap.String("homeDir", homeDir),
			zap.Error(err),
		)

		return
	}

	res, err = f.Readdirnames(-1)
	f.Close()

	return
}

// GetCertificate locates file of certificate
func GetCertificate() (loc string) {
	res, _ := readDir(homeDir)

	for _, item := range res {
		fileName := item

		if !strings.HasPrefix(fileName, hostnames[0]) {
			continue
		}

		if strings.HasSuffix(fileName, "-key.pem") {
			continue
		}

		if strings.HasSuffix(fileName, ".pem") {
			loc = filepath.Join(homeDir, fileName)

			logger.Info("identified certificate",
				zap.String("certificate", loc),
			)

			return
		}
	}
	return
}

// GetPrivateKey locates file of private key
func GetPrivateKey() (loc string) {
	res, _ := readDir(homeDir)

	for _, item := range res {
		fileName := item

		if !strings.HasPrefix(fileName, hostnames[0]) {
			continue
		}

		if strings.HasSuffix(fileName, "-key.pem") {
			loc = filepath.Join(homeDir, fileName)

			logger.Info("identified private key",
				zap.String("key", loc),
			)
			return
		}
	}
	return
}

// AutoInjectCA generates and inserts Root CA into trust store
func AutoInjectCA() {
	args := []string{
		"-install",
	}

	cmd := exec.Command(execFile, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Start()
	err := cmd.Wait()
	stdoutS := stdout.String()
	stderrS := stderr.String()

	log := logger.WithOptions(zap.Fields(
		zap.String("binary", execFile),
		zap.String("stdout", stdoutS),
		zap.String("stderr", stderrS),
	))

	if err != nil {
		log.Warn("error in installing root ca")
	}
	log.Info("installed root ca")
}

// Generate certificate and key if they don't exist
func Generate() (certFile, keyFile string) {
	shouldWork := false
	cert := GetCertificate()
	key := GetPrivateKey()
	if cert == "" || key == "" {
		shouldWork = true
	}
	if !shouldWork {
		return
	}

	// execute binary
	args := make([]string, 0)
	args = append(args, hostnames...)

	cmd := exec.Command(execFile, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Start()
	err := cmd.Wait()
	stdoutS := stdout.String()
	stderrS := stderr.String()

	log := logger.WithOptions(zap.Fields(
		zap.String("binary", execFile),
	))

	// parse output in order to extract cert and key
	parts := make([]string, 0)
	tmp := strings.Split(stderrS, "\n")

	for _, row := range tmp {
		if strings.Contains(row, "-key.pem") {
			parts = strings.Split(row, " ")
			break
		}
	}
	for _, part := range parts {
		part = strings.Trim(part, "\"")
		part = strings.ReplaceAll(part, "./", "")
		it := 0

		if strings.HasSuffix(part, ".pem") {
			it++
			newLocation := filepath.Join(homeDir, part)
			os.Rename(part, newLocation)

			if it == 1 {
				certFile = newLocation
			} else {
				keyFile = newLocation
			}

			log.Info("moving file",
				zap.String("oldLocation", part),
				zap.String("newLocation", newLocation),
			)
		}
	}

	if err != nil {
		log.Warn("error in generating certificate and key",
			zap.String("stdout", stdoutS),
			zap.String("stderr", stderrS),
		)
	}
	log.Info("generated certificate and key",
		zap.String("stdout", stdoutS),
		zap.String("stderr", stderrS),
	)

	return
}

func getPreferredBinary() (preferredBinary string) {
	switch runtime.GOOS {
	case "linux":
		preferredBinary = "mkcert-linux"
	case "darwin":
		preferredBinary = "mkcert-darwin"
	default:
		preferredBinary = "mkcert-linux"
	}

	return
}

// Bootstrap prepares for cert injection
func Bootstrap() {
	// autoselect binary
	preferredBinary := getPreferredBinary()
	filePaths := []string{
		fmt.Sprintf("/bin/%s", preferredBinary),
		fmt.Sprintf("./backend/server/ca/binaries/%s", preferredBinary),
	}

	for _, filePath := range filePaths {
		execFile, _ = filepath.Abs(filePath)
		if libIO.Exists(execFile) {
			break
		}
		execFile = ""
	}

	// fallback to system binary
	if !libIO.Exists(execFile) {
		file, err := exec.LookPath("mkcert")
		execFile = file
		if err != nil {
			execFile = "mkcert"
		}
	}

	// checking local home dir
	usr, err := user.Current()
	if err != nil {
		logger.Warn(err.Error())
		homeDir = "/"
	} else {
		homeDir = usr.HomeDir
	}
}

// Setup generates RootCA and inserts it in trust store, generates server certificate and private key
func Setup(log *zap.Logger) {
	logger = log
	Bootstrap()
	Generate()

	AutoInjectCA()
}
