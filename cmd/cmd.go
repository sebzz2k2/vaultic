package cmd

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	storage "github.com/sebzz2k2/vaultic/kv_store/wal"
	"github.com/sebzz2k2/vaultic/lexer"
	"github.com/sebzz2k2/vaultic/utils"
)

var processors = map[lexer.TokenKind]any{
	lexer.CMD_GET:    get,
	lexer.CMD_SET:    set,
	lexer.CMD_DEL:    del,
	lexer.CMD_EXISTS: exists,
	lexer.CMD_KEYS:   keys,
}

func validateArgsAndCount(t []lexer.Token) (bool, error) {
	if len(t) == 0 {
		return false, fmt.Errorf("no tokens")
	}
	if utils.CmdArgs[strings.ToUpper(t[0].Value)] != len(t)-1 {
		return false, fmt.Errorf("wrong arg count")
	}
	for _, tok := range t[1:] {
		if tok.Kind != lexer.VALUE {
			return false, fmt.Errorf("invalid token %s", tok.Value)
		}
	}
	return true, nil
}

func ProcessCommand(tokens []lexer.Token) (string, error) {
	if len(tokens) == 0 {
		return "", fmt.Errorf("no tokens provided")
	}

	cmd := tokens[0]

	if _, ok := processors[cmd.Kind]; !ok {
		return "", fmt.Errorf("invalid command: %s", cmd.Value)
	}

	isValidArgCount, err := validateArgsAndCount(tokens)
	if err != nil {
		return "", err
	}
	if !isValidArgCount {
		return "", fmt.Errorf("invalid argument count for command: %s", cmd.Value)
	}

	fn, ok := processors[cmd.Kind]
	if !ok {
		return "", fmt.Errorf("unknown command: %s", cmd.Value)
	}

	fnValue := reflect.ValueOf(fn)

	args := tokens[1:]
	reflectArgs := make([]reflect.Value, len(args))
	for i, tok := range args {
		reflectArgs[i] = reflect.ValueOf(tok.Value)
	}

	results := fnValue.Call(reflectArgs)

	if len(results) != 2 {
		return "", fmt.Errorf("unexpected number of return values")
	}

	strResult, ok := results[0].Interface().(string)
	if !ok {
		return "", fmt.Errorf("first return value not string")
	}

	if !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	return strResult, err
}

func get(key string) (string, error) {
	file, err := os.Open(utils.FILENAME)
	if err != nil {
		return "", err
	}
	defer file.Close()
	start, end, bool := utils.GetIndexVal(key)
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

func set(key, val string) (string, error) {
	now := time.Now()
	epochSeconds := now.Unix()

	setVal, totalLen := storage.EncodeData(1, false, uint64(epochSeconds), false, key, val)

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

func del(key string) (string, error) {
	_, _, bool := utils.GetIndexVal(key)
	if !bool {
		return "(nil)", nil
	}
	now := time.Now()
	epochSeconds := now.Unix()

	setVal, _ := storage.EncodeData(1, true, uint64(epochSeconds), false, key, "(nil)")
	file, err := os.OpenFile(utils.FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(setVal)
	if err != nil {
		return "", err
	}
	utils.DeleteIndexKey(key)
	return "OK", nil
}

func exists(key string) (string, error) {
	_, _, bool := utils.GetIndexVal(key)
	if bool {
		return "true", nil
	}
	return "false", nil
}

func keys() (string, error) {
	keys := utils.GetAllKeys()
	if len(keys) == 0 {
		return "(nil)", nil
	}
	return strings.Join(keys, ", "), nil
}
