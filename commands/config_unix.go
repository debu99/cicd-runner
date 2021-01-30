// +build linux darwin freebsd openbsd

package commands

import (
	"os"
	"path/filepath"

	"github.com/debu99/cicd-runner/helpers"
)

func getDefaultConfigDirectory() string {
	if os.Getuid() == 0 {
		return "/etc/cicd-runner"
	} else if homeDir := helpers.GetHomeDir(); homeDir != "" {
		return filepath.Join(homeDir, ".cicd-runner")
	} else if currentDir := helpers.GetCurrentWorkingDirectory(); currentDir != "" {
		return currentDir
	}
	panic("Cannot get default config file location")
}
