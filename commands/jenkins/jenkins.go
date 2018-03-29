package jenkins

import (
	"github.com/bndr/gojenkins"
	"github.com/jburandt/gobot/bot"

	//"flag"
	"os"
)

func Jenkins(command *bot.Cmd) (string, error) {
	// subcommands
	//listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	//deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)

	// list
	//listBuilds := listCommand.String("builds", "", "List job name to show builds")
	jenkins := gojenkins.CreateJenkins(nil, os.Getenv("JENKINS_URL"), os.Getenv("JENKINS_EMAIL"), os.Getenv("JENKINS_TOKEN"))
	_, err := jenkins.Init()

	if err != nil {
		bot.Reply(command.Channel, "Something Went Wrong: " + err.Error(), command.User)
	}

	/* GOOD CODE
		build, err := jenkins.GetJob("sensu-rancher-remediator/job/sensu-rancher-remediator/job/master")
		if err != nil {
			panic("Job Does Not Exist")
		}

		lastSuccessBuild, err := build.GetLastSuccessfulBuild()
		if err != nil {
			panic("Last SuccessBuild does not exist")
		}

		//duration := lastSuccessBuild.GetDuration()
		fmt.Println(lastSuccessBuild)
	*/
    jobInput := command.Args[0]
	jobName := jobInput+"/job/"+jobInput+"/job/master"
	builds, err := jenkins.GetAllBuildIds(jobName)

	if err != nil {
		bot.Reply(command.Channel, "Could not find job with that name. Error number: " + err.Error(), command.User)
	}

	for _, build := range builds {
		bot.Reply(command.Channel, "Found build from master branch:"+build.URL, command.User)
	}
	return "", nil
}


