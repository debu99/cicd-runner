package commands

import (
	"github.com/debu99/cicd-runner/helpers"
)

func getDefaultConfigDirectory() string {
	if currentDir := helpers.GetCurrentWorkingDirectory(); currentDir != "" {
		return currentDir
	}

	panic("Cannot get default config file location")
}
