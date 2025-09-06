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
		return false, fmt.Errorf("No enough tokens provided")
	}
	if utils.CmdArgs[strings.ToUpper(t[0].Value)] != len(t)-1 {
		return false, fmt.Errorf("Wrong argument count for command: %s", t[0].Value)
	}
	for _, tok := range t[1:] {
		if tok.Kind != lexer.VALUE {
			return false, fmt.Errorf("Invalid token: %s", tok.Value)
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
		return "", fmt.Errorf("No enough tokens provided")
	}

	cmd := tokens[0]

	if _, ok := processors[cmd.Kind]; !ok {
		return "", fmt.Errorf("Invalid command: %s", cmd.Value)
	}

	isValidArgCount, err := validateArgsAndCount(tokens)
	if err != nil {
		return "", err
	}
	if !isValidArgCount {
		return "", fmt.Errorf("Wrong argument count for command: %s", cmd.Value)
	}

	fn, ok := processors[cmd.Kind]
	if !ok {
		return "", fmt.Errorf("Unknown command: %s", cmd.Value)
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
		return "", fmt.Errorf("Unexpected return values")
	}

	strResult, ok := results[0].Interface().(string)
	if !ok {
		return "", fmt.Errorf("First return value is not a string")
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
			return "(nil)", nil
		}
		return "", fmt.Errorf("Failed to open WAL file")
	}
	defer file.Close()
	start, end, found := p.idx.Get(key)
	if found {
		_, err := file.Seek(int64(start), io.SeekStart)
		if err != nil {
			return "", fmt.Errorf("Failed to seek in WAL file")
		}
		b := make([]byte, end-start)
		n, err := file.Read(b)
		if err != nil {
			return "", fmt.Errorf("Failed to read from WAL file")
		}
		return string(b[:n]), nil
	}
	return "(nil)", nil
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
	return "OK", nil
}

func (p *Protocol) del(key string) (string, error) {
	_, _, found := p.idx.Get(key)
	if !found {
		return "(nil)", nil
	}
	now := time.Now()
	epochSeconds := now.Unix()

	setVal, _ := p.wal.EncodeWAL(1, true, uint64(epochSeconds), false, key, "(nil)")
	file, err := os.OpenFile(utils.FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to open WAL file")
	}

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return "", fmt.Errorf("Failed to seek in WAL file")
	}
	defer file.Close()
	_, err = file.Write(setVal)
	if err != nil {
		return "", fmt.Errorf("Failed to write to WAL file")
	}
	p.idx.Del(key)
	return "OK", nil
}

func (p *Protocol) exists(key string) (string, error) {
	_, _, found := p.idx.Get(key)
	if found {
		return "true", nil
	}
	return "false", nil
}

func (p *Protocol) keys() (string, error) {
	keys := p.idx.Keys()
	if len(keys) == 0 {
		return "(nil)", nil
	}
	return strings.Join(keys, ", "), nil
}
