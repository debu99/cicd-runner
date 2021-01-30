package main

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/debu99/cicd-runner/common"
	cli_helpers "github.com/debu99/cicd-runner/helpers/cli"
	"github.com/debu99/cicd-runner/log"

	_ "github.com/debu99/cicd-runner/cache/azure"
	_ "github.com/debu99/cicd-runner/cache/gcs"
	_ "github.com/debu99/cicd-runner/cache/s3"
	_ "github.com/debu99/cicd-runner/commands"
	_ "github.com/debu99/cicd-runner/commands/helpers"
	_ "github.com/debu99/cicd-runner/executors/custom"
	_ "github.com/debu99/cicd-runner/executors/docker"
	_ "github.com/debu99/cicd-runner/executors/docker/machine"
	_ "github.com/debu99/cicd-runner/executors/kubernetes"
	_ "github.com/debu99/cicd-runner/executors/parallels"
	_ "github.com/debu99/cicd-runner/executors/shell"
	_ "github.com/debu99/cicd-runner/executors/ssh"
	_ "github.com/debu99/cicd-runner/executors/virtualbox"
	_ "github.com/debu99/cicd-runner/helpers/secrets/resolvers/vault"
	_ "github.com/debu99/cicd-runner/shells"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// log panics forces exit
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "a CICD Runner"
	app.Version = common.AppVersion.ShortLine()
	cli.VersionPrinter = common.AppVersion.Printer
	app.Authors = []cli.Author{
		{
			Name:  "Clix.dev",
			Email: "support@Clix.dev",
		},
	}
	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logrus.Fatalln("Command", command, "not found.")
	}

	cli_helpers.InitCli()
	cli_helpers.LogRuntimePlatform(app)
	cli_helpers.SetupCPUProfile(app)
	cli_helpers.FixHOME(app)
	cli_helpers.WarnOnBool(os.Args)

	log.ConfigureLogging(app)

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
