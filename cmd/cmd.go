package cmd

import (
	"io"
	"os"

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
	offset, bool := utils.GetIndexVal(kn[0])
	if bool {
		_, err = file.Seek(int64(offset), 0)
		if err != nil {
			return "", err
		}
		var result []byte
		buf := make([]byte, 1)
		for {
			_, err := file.Read(buf)
			if err != nil {
				break
			}
			if buf[0] == utils.NEWLINE[0] {
				break
			}
			result = append(result, buf[0])
		}
		return string(result), nil
	}
	return "", nil
}

func (g GET) Validate(argsCount int) bool {
	return utils.CmdArgs[utils.CommandGet] == argsCount
}

type SET struct{}

func (s SET) Process(kv []string) (string, error) {
	key := kv[0]
	val := kv[1]

	setVal := key + utils.DELIMITER + val + utils.NEWLINE
	file, err := os.OpenFile(utils.FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	size, err := file.Seek(0, io.SeekEnd)
	offset := int(size) + len(key) + 1
	utils.SetIndexKey(key, offset)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.WriteString(setVal)
	if err != nil {
		return "", err
	}
	return utils.SUCCESS, nil
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
