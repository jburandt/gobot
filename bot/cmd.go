package bot

import (
	"fmt"
	"log"
	"strings"
)

// Struct that holds all of the commands
type Cmd struct {
	Raw         string       // Raw is full string passed to the command
	Channel     string       // Channel where the command was called
	ChannelData *ChannelData // More info about the channel, including network
	User        *User        // User who sent the message
	Message     string       // Full string without the prefix
	MessageData *Message     // Message with extra flags
	Command     string       // Command is the first argument passed to the bot
	RawArgs     string       // Raw arguments after the command
	Args        []string     // Arguments as array
}

type ChannelData struct {
	Protocol  string // The chat protocol (slack)
	Server    string // The server hostname the message was sent on
	Channel   string // The channel name the message appeared in
	IsPrivate bool   // Whether the channel is a group or private chat
}

type User struct {
	ID       string
	Nick     string
	RealName string
	IsBot    bool
}

type Message struct {
	Text     string // The actual content of this Message
	IsAction bool   // Bot marks true if triggered @bot do something
}

// Lets you register your own command
type customCommand struct {
	Version     int // Set version v1 or v2
	Cmd         string // The command being passed
	CmdFuncV1   CmdFunc // Creates type of *cmd and returns via string
	CmdFuncV2   CmdFuncChan // Creates type of *cmd and returns via channel
	Description string // Command description for help
	ExampleArgs string // help examples
}

// CmdReturn is the return result message of V1 commands
type CmdReturn struct {
	Channel string // The channel where the bot should send the message
	Message string // The message to be sent
}

// CmdReturnChan is the return result message of V2 commands
type CmdReturnChan struct {
	Channel string
	Message chan string
	Done    chan bool
}

const (
	v1 = iota
	v2
	helpDescripton    = "Command description: %s"
	helpUsage         = "Example: %s%s %s"
	availableCommands = "Available commands: `%v`"
	helpAboutCommand  = "Type: *%s help command* to see details about a specific command."
	helpCommand       = "help"
	commandNotAvailable   = "Command %v not available."
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
)

type CmdFunc func(cmd *Cmd) (string, error)
type CmdFuncChan func(cmd *Cmd) (CmdReturnChan, error)

var commands = make(map[string]*customCommand)

// Adds a command that responds in slack via string
func AddCommand(command, description, exampleArgs string, cmdFunc CmdFunc) {
	commands[command] = &customCommand{
		Version:     v1,
		Cmd:         command,
		CmdFuncV1:   cmdFunc,
		Description: description,
		ExampleArgs: exampleArgs,
	}
}

// Adds a command that responds in slack via a channel
func AddCommandChannel(command, description, exampleArgs string, cmdFunc CmdFuncChan) {
	commands[command] = &customCommand{
		Version:     v2,
		Cmd:         command,
		CmdFuncV2:   cmdFunc,
		Description: description,
		ExampleArgs: exampleArgs,
	}
}

func (b *Bot) handleCmd(c *Cmd) {
	cmd := commands[c.Command]

	if cmd == nil {
		log.Printf("Command not found %v", c.Command)
		return
	}

	switch cmd.Version {
	case v1:
		message, err := cmd.CmdFuncV1(c)
		b.checkCmdError(err, c)
		if message != "" {
			b.handlers.Response(c.Channel, message, c.User)
		}
	case v2:
		result, err := cmd.CmdFuncV2(c)
		b.checkCmdError(err, c)
		if result.Channel == "" {
			result.Channel = c.Channel
		}
		for {
			select {
			case message := <-result.Message:
				if message != "" {
					b.handlers.Response(result.Channel, message, c.User)
				}
			case <-result.Done:
				return
			}
		}
	}
}

func (b *Bot) checkCmdError(err error, c *Cmd) {
	if err != nil {
		errorMsg := fmt.Sprintf(errorExecutingCommand, c.Command, err.Error())
		log.Printf(errorMsg)
		b.handlers.Response(c.Channel, errorMsg, c.User)
	}
}

func (b *Bot) help(c *Cmd) {
	cmd, _ := parse("<@" + BotUserID + ">"+c.RawArgs, c.ChannelData, c.User)
	if cmd == nil {
		b.showRegisteredCommands(c.Channel, c.User)
		return
	}

	command := commands[cmd.Command]
	if command == nil {
		b.showRegisteredCommands(c.Channel, c.User)
		return
	}

	b.showHelp(cmd, command)
}

func (b *Bot) showHelp(c *Cmd, help *customCommand) {
	if help.Description != "" {
		b.handlers.Response(c.Channel, fmt.Sprintf(helpDescripton, help.Description), c.User)
	}
	b.handlers.Response(c.Channel, fmt.Sprintf(helpUsage, "<@" + BotUserID + "> ", c.Command, help.ExampleArgs), c.User)
}

func (b *Bot) showRegisteredCommands(channel string, sender *User) {
	var cmds []string
	for k := range commands {
		cmds = append(cmds, k)
	}
	b.handlers.Response(channel, fmt.Sprintf(helpAboutCommand, "<@" + BotUserID + ">"), sender)
	b.handlers.Response(channel, fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")), sender)
}
