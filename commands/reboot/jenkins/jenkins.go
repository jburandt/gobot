package jenkins

import (
	"github.com/bndr/gojenkins"
	"github.com/jburandt/gobot/bot"
	//"fmt"
)

func Jenkins(command *bot.Cmd) (string, error) {
	jenkins := gojenkins.CreateJenkins(nil, "<URL>", "<EMAIL>", "<API KEY>")
	_, err := jenkins.Init()

	if err != nil {
		panic("Something Went Wrong")
	}

	jobs, err := jenkins.GetAllJobs()
	if err != nil {
		panic(err)
	}
	for _, job := range jobs {
		bot.Reply(command.Channel, "Found job: "+job.Base, command.User)
	}

	return "", nil
}
