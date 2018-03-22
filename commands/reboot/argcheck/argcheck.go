package argcheck

import (
	"fmt"
	"github.com/jburandt/gobot/bot"
)

func Argcheck(command *bot.Cmd) (string, error) {
	one := fmt.Sprintf("Here are your raw args: %s", command.RawArgs)
	two := fmt.Sprintf("Arg0 %s", command.Args[0])
	three := fmt.Sprintf("Arg1 %s", command.Args[1])
	four := fmt.Sprintf("Arg2 %s", command.Args[2])
	//fmt.Sprintf("Arg3 %s", command.Args[3])
	return fmt.Sprintf("%s \n %s \n %s \n %s", one, two, three, four), nil
}
