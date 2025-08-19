package protocol

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sebzz2k2/vaultic/internal/index"
	"github.com/sebzz2k2/vaultic/internal/protocol/lexer"
	"github.com/sebzz2k2/vaultic/internal/wal"

	"github.com/sebzz2k2/vaultic/pkg/config"
	"github.com/sebzz2k2/vaultic/pkg/utils"
)

type Protocol struct {
	idx *index.Index
	wal *wal.WAL
}

var processors = map[lexer.TokenKind]any{
	lexer.CMD_GET:    (*Protocol).get,
	lexer.CMD_SET:    (*Protocol).set,
	lexer.CMD_DEL:    (*Protocol).del,
	lexer.CMD_EXISTS: (*Protocol).exists,
	lexer.CMD_KEYS:   (*Protocol).keys,
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

func NewProtocol(wal *wal.WAL) *Protocol {
	return &Protocol{
		//TODO remove hardcoded index file name
		idx: index.NewIndex("vaultic", wal),
		wal: wal,
	}
}

func (p *Protocol) ProcessCommand(tokens []lexer.Token) (string, error) {
	fmt.Println("Processing command:", tokens)
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
	reflectArgs := make([]reflect.Value, len(args)+1)
	reflectArgs[0] = reflect.ValueOf(p)
	for i, tok := range args {
		reflectArgs[i+1] = reflect.ValueOf(tok.Value)
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

func (p *Protocol) get(key string) (string, error) {
	file, err := os.Open(utils.FILENAME)
	if err != nil {
		if os.IsNotExist(err) {
			return config.NilMessage, nil
		}
		return "", fmt.Errorf(config.ErrorFileOpen)
	}
	defer file.Close()
	start, end, found := p.idx.Get(key)
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

func (p *Protocol) set(key, val string) (string, error) {
	now := time.Now()
	epochSeconds := now.Unix()

	setVal, totalLen := p.wal.EncodeWAL(1, false, uint64(epochSeconds), false, key, val)

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
	p.idx.Set(key, start, offset)
	defer file.Close()
	_, err = file.Write(setVal)
	if err != nil {
		return "", err
	}
	return config.OKMessage, nil
}

func (p *Protocol) del(key string) (string, error) {
	_, _, found := p.idx.Get(key)
	if !found {
		return config.NilMessage, nil
	}
	now := time.Now()
	epochSeconds := now.Unix()

	setVal, _ := p.wal.EncodeWAL(1, true, uint64(epochSeconds), false, key, config.NilMessage)
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
	p.idx.Del(key)
	return config.OKMessage, nil
}

func (p *Protocol) exists(key string) (string, error) {
	_, _, found := p.idx.Get(key)
	if found {
		return config.TrueMessage, nil
	}
	return config.FalseMessage, nil
}

func (p *Protocol) keys() (string, error) {
	keys := p.idx.Keys()
	if len(keys) == 0 {
		return config.NilMessage, nil
	}
	return strings.Join(keys, ", "), nil
}
