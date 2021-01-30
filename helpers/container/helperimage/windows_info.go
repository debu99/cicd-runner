package helperimage

import (
	"fmt"

	"github.com/debu99/cicd-runner/helpers/container/windows"
	"github.com/debu99/cicd-runner/shells"
)

const (
	baseImage1809 = "servercore1809"
	baseImage1903 = "servercore1903"
	baseImage1909 = "servercore1909"
	baseImage2004 = "servercore2004"

	windowsSupportedArchitecture = "x86_64"
)

var helperImages = map[string]string{
	windows.V1809: baseImage1809,
	windows.V1903: baseImage1903,
	windows.V1909: baseImage1909,
	windows.V2004: baseImage2004,
}

var powerShellCmdArgs = []string{
	"-NoProfile",
	"-NoLogo",
	"-InputFormat",
	"text",
	"-OutputFormat",
	"text",
	"-NonInteractive",
	"-ExecutionPolicy",
	"Bypass",
	"-Command",
	"-",
}

type windowsInfo struct{}

func (w *windowsInfo) Create(revision string, cfg Config) (Info, error) {
	baseImage, err := w.baseImage(cfg.OperatingSystem)
	if err != nil {
		return Info{}, fmt.Errorf("detecting base image: %w", err)
	}

	shell := cfg.Shell
	if cfg.Shell == "" {
		// TODO: Replace with shells.SNPwsh in 14.0 in https://gitlab.com/gitlab-org/gitlab-runner/-/issues/26419
		shell = shells.SNPowershell
	}
	return Info{
		Architecture:            windowsSupportedArchitecture,
		Name:                    imageName(cfg.GitLabRegistry),
		Tag:                     fmt.Sprintf("%s-%s-%s", windowsSupportedArchitecture, revision, baseImage),
		IsSupportingLocalImport: false,
		Cmd:                     append([]string{shell}, powerShellCmdArgs...),
	}, nil
}

func (w *windowsInfo) baseImage(operatingSystem string) (string, error) {
	version, err := windows.Version(operatingSystem)
	if err != nil {
		return "", err
	}

	baseImage, ok := helperImages[version]
	if !ok {
		return "", windows.NewUnsupportedWindowsVersionError(operatingSystem)
	}

	return baseImage, nil
}
