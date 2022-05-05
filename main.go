package main

import (
	"encoding/json"
	"fmt"
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
		log.Panic(err)
	}

	machine := config.Machines[0]

	auth_key, _ := homedir.Expand(config.AuthKey)
	auth, err := goph.Key(auth_key, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := goph.New(machine.User, machine.Host, auth)
	if err != nil {
		log.Panicln(err)
	}

	defer client.Close()

	out, err := client.Run(config.Commands["squeue"])

	if err != nil {
		log.Panic(err)
	}

	parsed, err := utils.ParseSqueue(string(out), machine)

	if err != nil {
		log.Panic(err)
	}

	bytes, err := json.Marshal(parsed)

	fmt.Println(string(bytes))

	return

}

func main() {

	var err error

	app := &cli.App{
		Name: "ssl-tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				Value:   "config.yml",
			},
		},
		Commands: [](*cli.Command){
			{
				Name:   "jobs",
				Usage:  "List all the running jobs.",
				Action: GetStatus,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"fmt", "f"},
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
