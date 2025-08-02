package cmd

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	storage "github.com/sebzz2k2/vaultic/kv_store"
	"github.com/sebzz2k2/vaultic/lexer"
	"github.com/sebzz2k2/vaultic/pkg/config"
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
		return false, fmt.Errorf(config.ErrorNoTokens)
	}
	if utils.CmdArgs[strings.ToUpper(t[0].Value)] != len(t)-1 {
		return false, fmt.Errorf(config.ErrorWrongArgCount, t[0].Value)
	}
	for _, tok := range t[1:] {
		if tok.Kind != lexer.VALUE {
			return false, fmt.Errorf(config.ErrorInvalidToken, tok.Value)
		}
	}
	return true, nil
}

func ProcessCommand(tokens []lexer.Token) (string, error) {
	if len(tokens) == 0 {
		return "", fmt.Errorf(config.ErrorNoTokens)
	}

	cmd := tokens[0]

	if _, ok := processors[cmd.Kind]; !ok {
		return "", fmt.Errorf(config.ErrorInvalidCommand, cmd.Value)
	}

	isValidArgCount, err := validateArgsAndCount(tokens)
	if err != nil {
		return "", err
	}
	if !isValidArgCount {
		return "", fmt.Errorf(config.ErrorWrongArgCount, cmd.Value)
	}

	fn, ok := processors[cmd.Kind]
	if !ok {
		return "", fmt.Errorf(config.ErrorUnknownCommand, cmd.Value)
	}

	fnValue := reflect.ValueOf(fn)

	args := tokens[1:]
	reflectArgs := make([]reflect.Value, len(args))
	for i, tok := range args {
		reflectArgs[i] = reflect.ValueOf(tok.Value)
	}

	results := fnValue.Call(reflectArgs)

	if len(results) != 2 {
		return "", fmt.Errorf(config.ErrorUnexpectedReturn)
	}

	strResult, ok := results[0].Interface().(string)
	if !ok {
		return "", fmt.Errorf(config.ErrorFirstReturnNotStr)
	}

	if !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	return strResult, err
}

func get(key string) (string, error) {
	file, err := os.Open(utils.FILENAME)
	if err != nil {
		return "", fmt.Errorf(config.ErrorFileOpen)
	}
	defer file.Close()
	start, end, found := utils.GetIndexVal(key)
	if found {
		_, err := file.Seek(int64(start), io.SeekStart)
		if err != nil {
			return "", fmt.Errorf(config.ErrorFileSeek)
		}
		b := make([]byte, end-start)
		n, err := file.Read(b)
		if err != nil {
			return "", fmt.Errorf(config.ErrorFileRead)
		}
		return string(b[:n]), nil
	}
	return config.NilMessage, nil
}

func set(key, val string) (string, error) {
	now := time.Now()
	epochSeconds := now.Unix()

	setVal, totalLen := storage.EncodeWAL(1, false, uint64(epochSeconds), false, key, val)

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
	_, _, found := utils.GetIndexVal(key)
	if !found {
		return config.NilMessage, nil
	}
	now := time.Now()
	epochSeconds := now.Unix()

	setVal, _ := storage.EncodeWAL(1, true, uint64(epochSeconds), false, key, config.NilMessage)
	file, err := os.OpenFile(utils.FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf(config.ErrorFileOpen)
	}

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return "", fmt.Errorf(config.ErrorFileSeek)
	}
	defer file.Close()
	_, err = file.Write(setVal)
	if err != nil {
		return "", fmt.Errorf(config.ErrorFileWrite)
	}
	utils.DeleteIndexKey(key)
	return config.OKMessage, nil
}

func exists(key string) (string, error) {
	_, _, found := utils.GetIndexVal(key)
	if found {
		return config.TrueMessage, nil
	}
	return config.FalseMessage, nil
}

func keys() (string, error) {
	keys := utils.GetAllKeys()
	if len(keys) == 0 {
		return config.NilMessage, nil
	}
	return strings.Join(keys, ", "), nil
}
