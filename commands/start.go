package commands

import (
	"log"
	"runtime"
	"time"

	"github.com/hirakiuc/cloud-pubsub-worker-go/common"
	"github.com/hirakiuc/cloud-pubsub-worker-go/dispatcher"
	"github.com/urfave/cli"
)

// StartCommand describe the options of start command.
type StartCommand struct {
	dispatcher *dispatcher.Dispatcher
}

// Execute execute start command task.
func (c *StartCommand) Execute(context *cli.Context) error {
	log.Println("Start Command")

	config := dispatcher.Config{
		NumOfWorkers: runtime.NumCPU(),
		JobRate:      time.Second / 2,
		JobBurst:     1,
	}
	var err error
	c.dispatcher, err = dispatcher.NewDispatcher(&config)
	if err != nil {
		return err
	}

	return c.dispatcher.Start()
}

func init() {
	cmd := &StartCommand{}

	command := cli.Command{
		Name:   "start",
		Usage:  "Start dispatcher",
		Action: cmd.Execute,
	}

	common.RegisterCommand(command)
}
