package access

import (
	"errors"
	"os/exec"
	"github.com/jburandt/gobot/bot"

)

var (
	allowedUsers = map[string]bool{"jburandt": true}
	errAllowedUsers = errors.New("User does not have permission to execute this command")
	)

func cmd(command *bot.Cmd) (string, error) {
	if _, check := allowedUsers[command.Args[0]]; check {
		return "", errAllowedUsers
	}
	cmd := exec.Command("/bin/bash", "-c", command.RawArgs)
	data, err := cmd.CombinedOutput()
	return string(data), err
}
