package debug

import (
	"fmt"

	"github.com/jburandt/gobot/bot"
)

func Debug(command *bot.Cmd) (msg string, err error) {
    channeldata := command.ChannelData.Channel
    channelserver := command.ChannelData.Server
    isprivate := fmt.Sprintf("%t", command.ChannelData.IsPrivate )
    userid := command.User.ID
    nick := command.User.Nick
    realname := command.User.RealName
    botid := bot.BotUserID
	bot.Reply(command.Channel, "##### PRINTING DEBUG INFO #####", command.User)
	bot.Reply(command.Channel, "Your real name is listed as: " +  realname, command.User)
	bot.Reply(command.Channel, "Your channel info is: " +  channeldata, command.User)
	bot.Reply(command.Channel, "Your channel server is: " +  channelserver, command.User)
	bot.Reply(command.Channel, "This channel is private: " +  isprivate, command.User)
	bot.Reply(command.Channel, "Your user id is listed as: " +  userid, command.User)
	bot.Reply(command.Channel, "Your nick name is listed as : " +  nick, command.User)
	bot.Reply(command.Channel, "Your bot id is listed as: " +  botid, command.User)
	return "", err
}
