package cmd

import (
	"fmt"

	"github.com/sebzz2k2/vaultic/utils"
)

type Command interface {
	Process()
	Validate(argsCount int) bool
}

type GET struct {
}

func (g GET) Process() {
	fmt.Println("GET")
}
func (g GET) Validate(argsCount int) bool {
	return utils.CmdArgs[utils.CommandGet] == argsCount
}

type SET struct{}

func (s SET) Process() {
	fmt.Println("SET")
}
func (s SET) Validate(argsCount int) bool {
	return utils.CmdArgs[utils.CommandGet] == argsCount
}

func CommandFactory(command string) Command {
	switch command {
	case utils.CommandGet:
		{
			return GET{}
		}
	case utils.CommandSet:
		{
			return SET{}
		}

	default:
		return nil
	}
}
