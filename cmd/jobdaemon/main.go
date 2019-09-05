package main

import (
	"os"
	"path"

	"github.com/hirakiuc/jobworker-go/common"
	"go.uber.org/zap"

	"github.com/urfave/cli"

	_ "github.com/hirakiuc/jobworker-go/commands"
)

func main() {
	common.ConfigureLogger()
	logger := common.GetLogger()
	defer logger.Sync() // flushes buffer, if any

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "a SampleJob Runner"
	app.Version = "1.0"
	app.Authors = []cli.Author{
		{
			Name:  "hirakiuc",
			Email: "hirakiuc@gmail.com",
		},
	}
	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logger.Fatal("Command not found.",
			zap.String("name", command),
		)
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("Command failed",
			zap.Error(err),
		)
	}
}
