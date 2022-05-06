package main

import (
	"log"
	"os"

	"github.com/chazeon/hpc-cli/utils"
	"github.com/melbahja/goph"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

func GetStatus(c *cli.Context) (err error) {

	config, err := utils.LoadConfig(c.String("config"))

	if err != nil {
		return err
	}

	machine := config.Machines[0]

	auth_key, _ := homedir.Expand(config.AuthKey)
	auth, err := goph.Key(auth_key, "")

	if err != nil {
		return err
	}

	client, err := goph.New(machine.User, machine.Host, auth)

	if err != nil {
		return err
	}

	defer client.Close()

	out, err := client.Run(config.Commands["squeue"])

	if err != nil {
		return err
	}

	jobs, err := utils.ParseJobs(string(out), machine)

	if err != nil {
		return err
	}

	utils.ShowJobs(jobs, c.String("format"))

	return

}

func main() {

	var err error

	app := &cli.App{
		Name: "hpc-cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				Value:   "config.yml",
			},
			&cli.StringSliceFlag{
				Name:    "machine",
				Aliases: []string{"m"},
				Usage:   "Specify the `MACHINE` to apply on.",
				Value:   nil,
			},
		},
		Commands: [](*cli.Command){
			{
				Name:   "exec",
				Usage:  "Execute the command.",
				Action: utils.ExecCommand,
			},
			{
				Name:   "jobs",
				Usage:  "List all the running jobs.",
				Action: GetStatus,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"fmt", "f"},
						Value:   "table",
					},
				},
			},
		},
	}

	err = app.Run(os.Args)

	if err != nil {
		log.Panic(err)
	}

	return
}
