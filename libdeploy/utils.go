package libdeploy

import (
	"os"
	"path/filepath"

	pathutils "github.com/xozrc/deploy/pkg/pathutils"
)

var (
	baseDir = os.Getenv("MACHINE_STORAGE_PATH")

	defaultDeployPath   = ".deploy"
	defaultMachinesPath = "machines"
	defaultCertPath     = "certs"
)

func GetBaseDir() string {
	if baseDir == "" {
		baseDir = filepath.Join(pathutils.GetHomeDir(), defaultDeployPath)
	}
	return baseDir
}

func GetMachineDir() string {
	return filepath.Join(GetBaseDir(), defaultMachinesPath)
}

func GetMachineCertDir() string {
	return filepath.Join(GetBaseDir(), defaultCertPath)
}
