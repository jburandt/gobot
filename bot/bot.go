package bot

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/robfig/cron"
	"regexp"
	"strings"
	"github.com/mattn/go-shellwords"
	"errors"
)

type Bot struct {
	handlers     *Handlers
	cron         *cron.Cron
	disabledCmds []string
}

// MessageHandler must be implemented by the protocol to handle the bot responses
type MessageHandler func(target, message string, sender *User)

// Handlers that must be registered to receive callbacks from the bot
type Handlers struct {
	Response MessageHandler
}

// New configures a new bot instance
func New(h *Handlers) *Bot {
	b := &Bot{
		handlers: h,
		cron:     cron.New(),
	}
	return b
}

var (
	re = regexp.MustCompile("\\s+") // Matches one or more spaces
)

func parse(s string, channel *ChannelData, user *User) (*Cmd, error) {

	c := &Cmd{Raw: s}
	s = strings.TrimSpace(s)

	if !strings.HasPrefix(s, "<@" + BotUserID + ">") {
		return nil, nil
	}

	c.Channel = strings.TrimSpace(channel.Channel)
	c.ChannelData = channel
	c.User = user

	// Trim the prefix and extra spaces
	c.Message = strings.TrimPrefix(s, "<@" + BotUserID + ">")
	c.Message = strings.TrimSpace(c.Message)

	// check if we have the command and not only the prefix
	if c.Message == "" {
		return nil, nil
	}

	// get the command
	pieces := strings.SplitN(c.Message, " ", 2)

	c.Command = pieces[0]

	if len(pieces) > 1 {
		// get the arguments and remove extra spaces
		c.RawArgs = strings.TrimSpace(pieces[1])
		parsedArgs, err := shellwords.Parse(c.RawArgs)
		if err != nil {
			return nil, errors.New("Error parsing arguments: " + err.Error())
		}
		c.Args = parsedArgs
	}

	c.MessageData = &Message{
		Text: c.Message,
	}

	return c, nil
}

// MessageReceived must be called by the protocol upon receiving a message
func (b *Bot) MessageReceived(channel *ChannelData, message *Message, sender *User) {
	command, err := parse(message.Text, channel, sender)

	if err != nil {
		b.handlers.Response(channel.Channel, err.Error(), sender)
		return
	}

	if command == nil {

		for i := 0; i <= 50; i++ {
			fmt.Println(message.Text)
			fmt.Printf("%#v\n", message)
			fmt.Println(channel.Channel)
			fmt.Println(sender)
		}
		return
	}

	switch command.Command {
	case helpCommand:
		b.help(command)
	default:
		b.handleCmd(command)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
