package docker

import (
	"github.com/debu99/cicd-runner/common"
	"github.com/debu99/cicd-runner/executors"
	"github.com/debu99/cicd-runner/executors/docker/internal/volumes/parser"
	"github.com/debu99/cicd-runner/executors/docker/internal/volumes/permission"
)

func init() {
	options := executors.ExecutorOptions{
		DefaultCustomBuildsDirEnabled: true,
		DefaultBuildsDir:              `c:\builds`,
		DefaultCacheDir:               `c:\cache`,
		SharedBuildsDir:               false,
		Shell: common.ShellScriptInfo{
			Shell:         "powershell",
			Type:          common.NormalShell,
			RunnerCommand: "cicd-runner-helper",
		},
		ShowHostname: true,
		Metadata: map[string]string{
			metadataOSType: osTypeWindows,
		},
	}

	creator := func() common.Executor {
		e := &commandExecutor{
			executor: executor{
				AbstractExecutor: executors.AbstractExecutor{
					ExecutorOptions: options,
				},
				volumeParser: parser.NewWindowsParser(),
			},
		}

		e.newVolumePermissionSetter = func() (permission.Setter, error) {
			return permission.NewDockerWindowsSetter(), nil
		}

		e.SetCurrentStage(common.ExecutorStageCreated)
		return e
	}

	featuresUpdater := func(features *common.FeaturesInfo) {
		features.Variables = true
		features.Image = true
		features.Services = true
		features.Session = false
		features.Terminal = false
	}

	common.RegisterExecutorProvider("docker-windows", executors.DefaultExecutorProvider{
		Creator:          creator,
		FeaturesUpdater:  featuresUpdater,
		DefaultShellName: options.Shell.Shell,
	})
}
