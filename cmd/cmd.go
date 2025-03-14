package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/sebzz2k2/vaultic/utils"
)

type Command interface {
	Process(kv []string) (string, error)
	Validate(argsCount int) bool
}

type GET struct {
}

const ()

func (g GET) Process(kn []string) (string, error) {
	file, err := os.Open(utils.FILENAME)
	if err != nil {
		return "", err
	}
	defer file.Close()
	pattern := fmt.Sprintf(`^%s(.*)`, kn[0]+utils.DELIMITER)
	re := regexp.MustCompile(pattern)
	var lastMatch string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); matches != nil {
			lastMatch = matches[1] // Extract content after "key:"
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return lastMatch, nil
}

func (g GET) Validate(argsCount int) bool {
	return utils.CmdArgs[utils.CommandGet] == argsCount
}

type SET struct{}

func (s SET) Process(kv []string) (string, error) {
	key := kv[0]
	val := kv[1]

	setVal := key + utils.DELIMITER + val + utils.NEWLINE

	return "", utils.WriteToFile(utils.FILENAME, setVal)
}
func (s SET) Validate(argsCount int) bool {
	return utils.CmdArgs[utils.CommandSet] == argsCount
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
