package utils

import (
	"strings"

	"github.com/melbahja/goph"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

func GetCommand(c *cli.Context) (cmd string) {

	args := make([]string, 0)

	for i := 0; i < c.NArg(); i++ {
		arg := c.Args().Get(i)
		args = append(args, string(arg))
	}

	cmd = strings.Join(args, " ")

	return cmd
}

func GetClients(machines []Machine, auth goph.Auth) (clients []*goph.Client, err error) {

	clients = make([]*goph.Client, 0)

	for _, machine := range machines {

		client, err := goph.New(machine.User, machine.Host, auth)

		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	return
}

type CommandResult struct {
	Machine Machine
	out     []byte
	err     error
}

func RunCommand(machine Machine, client *goph.Client, cmd string, ch chan CommandResult) {

	cr := CommandResult{
		Machine: machine,
	}

	cr.out, cr.err = client.Run(cmd)
	ch <- cr
}

func ExecCommand(c *cli.Context) (err error) {

	cmd := GetCommand(c)

	config, err := LoadConfig(c.String("config"))

	if err != nil {
		return err
	}

	auth_key, _ := homedir.Expand(config.AuthKey)
	auth, err := goph.Key(auth_key, "")

	if err != nil {
		return err
	}

	var machines []Machine

	if len(c.StringSlice("machine")) == 0 {

		machines = config.Machines

	} else {

		machines = []Machine{}
		mIndex := map[string]Machine{}

		for _, machine := range config.Machines {
			mIndex[machine.Name] = machine
		}

		for _, name := range c.StringSlice("machine") {
			if val, ok := mIndex[name]; ok {
				machines = append(machines, val)
			}
		}
	}

	clients, err := GetClients(machines, auth)

	if err != nil {
		return err
	}

	ch := make(chan CommandResult, len(clients))

	for i, machine := range machines {

		client := clients[i]
		go RunCommand(machine, client, cmd, ch)
		defer client.Close()

	}

	for range clients {

		result := <-ch

		if result.err != nil {
			return err
		}

		println(result.Machine.Host, ":")
		println(string(result.out))

	}

	return nil

}
