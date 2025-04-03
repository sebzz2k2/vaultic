package cmd

import (
	"io"
	"os"
	"time"

	"github.com/sebzz2k2/vaultic/kv_store/storage"
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
	start, end, bool := utils.GetIndexVal(kn[0])
	if bool {
		_, err := file.Seek(int64(start), io.SeekStart)
		if err != nil {
			return "", err
		}
		b := make([]byte, end-start)
		n, err := file.Read(b)
		if err != nil {
			return "", err
		}
		return string(b[:n]), nil
	}
	return "(nil)", nil
}

func (g GET) Validate(argsCount int) bool {
	return utils.CmdArgs[utils.CommandGet] == argsCount
}

type SET struct{}

func (s SET) Process(kv []string) (string, error) {
	key := kv[0]
	val := kv[1]
	now := time.Now()
	epochSeconds := now.Unix()

	setVal, totalLen := storage.EncodeData(1, false, uint64(epochSeconds), key, val)

	file, err := os.OpenFile(utils.FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return "", err
	}
	offset := uint32(size) + uint32(totalLen)
	start := offset - uint32(len(val))
	utils.SetIndexKey(key, start, offset)
	defer file.Close()
	_, err = file.Write(setVal)
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
