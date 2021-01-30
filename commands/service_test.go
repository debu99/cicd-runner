package commands

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"

	"github.com/debu99/cicd-runner/helpers"
)

func newTestGetServiceArgumentsCommand(t *testing.T, expectedArgs []string) func(*cli.Context) {
	return func(c *cli.Context) {
		arguments := getServiceArguments(c)

		for _, arg := range expectedArgs {
			assert.Contains(t, arguments, arg)
		}
	}
}

func testServiceCommandRun(command func(*cli.Context), args ...string) {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "test-command",
			Action: command,
			Flags:  getInstallFlags(),
		},
	}

	args = append([]string{"binary", "test-command"}, args...)
	_ = app.Run(args)
}

type getServiceArgumentsTestCase struct {
	cliFlags     []string
	expectedArgs []string
}

func TestGetServiceArguments(t *testing.T) {
	tests := []getServiceArgumentsTestCase{
		{
			expectedArgs: []string{
				"--working-directory", helpers.GetCurrentWorkingDirectory(),
				"--config", getDefaultConfigFile(),
				"--service", "cicd-runner",
				"--syslog",
			},
		},
		{
			cliFlags: []string{
				"--config", "/tmp/config.toml",
			},
			expectedArgs: []string{
				"--working-directory", helpers.GetCurrentWorkingDirectory(),
				"--config", "/tmp/config.toml",
				"--service", "cicd-runner",
				"--syslog",
			},
		},
		{
			cliFlags: []string{
				"--working-directory", "/tmp",
			},
			expectedArgs: []string{
				"--working-directory", "/tmp",
				"--config", getDefaultConfigFile(),
				"--service", "cicd-runner",
				"--syslog",
			},
		},
		{
			cliFlags: []string{
				"--service", "cicd-runner-service-name",
			},
			expectedArgs: []string{
				"--working-directory", helpers.GetCurrentWorkingDirectory(),
				"--config", getDefaultConfigFile(),
				"--service", "cicd-runner-service-name",
				"--syslog",
			},
		},
		{
			cliFlags: []string{
				"--syslog=true",
			},
			expectedArgs: []string{
				"--working-directory", helpers.GetCurrentWorkingDirectory(),
				"--config", getDefaultConfigFile(),
				"--service", "cicd-runner",
				"--syslog",
			},
		},
		{
			cliFlags: []string{
				"--syslog=false",
			},
			expectedArgs: []string{
				"--working-directory", helpers.GetCurrentWorkingDirectory(),
				"--config", getDefaultConfigFile(),
				"--service", "cicd-runner",
			},
		},
	}

	for id, testCase := range tests {
		t.Run(fmt.Sprintf("case-%d", id), func(t *testing.T) {
			testServiceCommandRun(newTestGetServiceArgumentsCommand(t, testCase.expectedArgs), testCase.cliFlags...)
		})
	}
}
