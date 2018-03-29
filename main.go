package main

import (
	"os"

	"github.com/jburandt/gobot/bot"
	"github.com/jburandt/gobot/commands/reboot"
	"github.com/jburandt/gobot/commands/debug"
	"github.com/jburandt/gobot/commands/argcheck"
	"github.com/jburandt/gobot/commands/jenkins"
	"github.com/jburandt/gobot/commands/copy"
)

func setup() {
	bot.AddCommand(
		"debug",
		"various debug commands to get slack info",
		"",
		debug.Debug)

	bot.AddCommandChannel(
		"reboot",
		"reboot any datacenter server",
		"",
		reboot.Reboot)

	bot.AddCommand(
		"argcheck",
		"gets count of args for testing",
		"",
		argcheck.Argcheck)

	bot.AddCommand(
		"jenkins",
		"hits jenkins api to show jobs",
		"",
		jenkins.Jenkins)
	bot.AddCommand(
		"copy",
		"hits jenkins api to show jobs",
		"",
		copy.Copy)

}

func main() {
	setup()
	bot.Run(os.Getenv("SLACK_TOKEN"))
}
