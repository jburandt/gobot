package reboot

import (
	"bufio"
	"bytes"

	"github.com/jburandt/gobot/bot"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
)

// to handle json keys
type rebootObject struct {
	Name    string
	Command string
}

var (
	fileNotFound = "config file not found"
	cmdNotFound  = "Error finding command argument"
)

// *bot.Cmd is the first arg to be passed to the bot
// bot.CmdReturnChan wants message from go routine and done = true
func Reboot(command *bot.Cmd) (result bot.CmdReturnChan, err error) {
	result = bot.CmdReturnChan{Message: make(chan string), Done: make(chan bool, 1)}

	// load json config file with names/commands
	filePath := "./commands/reboot/config.json"
	file, err1 := ioutil.ReadFile(filePath)
	if err1 != nil {
		bot.Reply(command.Channel, fileNotFound, command.User)
	}

	var scriptParse []rebootObject
	//userinput := "box4535345346" // faking user input
	userinput := command.Args[0] // real one
	err2 := json.Unmarshal(file, &scriptParse)
	if err2 != nil {
		fmt.Println("error:", err2)
		bot.Reply(command.Channel, fileNotFound, command.User)
	}
	//strip numbers off input to match json key
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}
	// loop through json file to find the match of user input to json name key
	cmdFound := false
	for k := range scriptParse {
		newinput := reg.ReplaceAllString(userinput, "")
		// keep running for loop until names do match
		if scriptParse[k].Name != newinput {
			continue
		}
		cmdFound = true

		// https://golang.org/ref/mem
		cmd := exec.Command("/bin/bash", "-c", scriptParse[k].Command)
		var b bytes.Buffer

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err.Error())
		}

		//Error
		cmd.Stderr = &b

		go func() {
			in := bufio.NewScanner(stdout)
			for in.Scan() {
				result.Message <- in.Text() // write each line to your log, or anything you need
			}
			result.Done <- true
		}()

		//Start command
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			cmd.Wait()
		}()
	}
	if cmdFound == false {
		result.Done <- true
		bot.Reply(command.Channel, cmdNotFound, command.User)
	}
	return result, nil
}
