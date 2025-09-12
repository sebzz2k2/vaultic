package resp

import (
	"bufio"
	"errors"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Decoder struct {
	reader *bufio.Reader
}

type RESPValue struct {
	Type   string
	String string
	Int    int64
	Float  float64
	BigInt *big.Int
	Bool   bool
	Array  []RESPValue
	Map    map[string]RESPValue
	Null   bool
	Error  string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{reader: bufio.NewReader(r)}
}

func (d *Decoder) Decode() (*RESPValue, error) {
	typeByte, err := d.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch typeByte {
	case respTypeToChar[SIMPLE_STRING]:
		return d.decodeSimpleString()
	case respTypeToChar[ERROR]:
		return d.decodeError()
	case respTypeToChar[INTEGER]:
		return d.decodeInteger()
	case respTypeToChar[BULK_STRING]:
		return d.decodeBulkString()
	case respTypeToChar[ARRAY]:
		return d.decodeArray()
	case respTypeToChar[NULL]:
		return d.decodeNull()
	case respTypeToChar[BOOLEAN]:
		return d.decodeBoolean()
	case respTypeToChar[DOUBLE]:
		return d.decodeDouble()
	case respTypeToChar[BIG_NUMBER]:
		return d.decodeBigNumber()
	case respTypeToChar[BULK_ERROR]:
		return d.decodeBulkError()
	case respTypeToChar[VERBATIM]:
		return d.decodeVerbatimString()
	case respTypeToChar[MAP]:
		return d.decodeMap()
	case respTypeToChar[ATTRIBUTES]:
		return d.decodeAttribute()
	case respTypeToChar[SET]:
		return d.decodeSet()
	case respTypeToChar[PUSH]:
		return d.decodePush()
	default:
		return nil, errors.New("unknown RESP type: " + string(typeByte))
	}
}

func (d *Decoder) readLine() (string, error) {
	line, err := d.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(line) < 2 || line[len(line)-2] != '\r' {
		return "", errors.New("invalid CRLF terminator")
	}
	return line[:len(line)-2], nil
}

func (d *Decoder) decodeSimpleString() (*RESPValue, error) {
	str, err := d.readLine()
	if err != nil {
		return nil, err
	}
	return &RESPValue{Type: SIMPLE_STRING, String: str}, nil
}

func (d *Decoder) decodeError() (*RESPValue, error) {
	str, err := d.readLine()
	if err != nil {
		return nil, err
	}
	return &RESPValue{Type: ERROR, Error: str}, nil
}

func (d *Decoder) decodeInteger() (*RESPValue, error) {
	str, err := d.readLine()
	if err != nil {
		return nil, err
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, err
	}
	return &RESPValue{Type: INTEGER, Int: val}, nil
}

func (d *Decoder) decodeBulkString() (*RESPValue, error) {
	lengthStr, err := d.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return &RESPValue{Type: BULK_STRING, Null: true}, nil
	}

	if length == 0 {
		d.readLine()
		return &RESPValue{Type: BULK_STRING, String: ""}, nil
	}

	data := make([]byte, length)
	_, err = io.ReadFull(d.reader, data)
	if err != nil {
		return nil, err
	}

	d.readLine()
	return &RESPValue{Type: BULK_STRING, String: string(data)}, nil
}

func (d *Decoder) decodeArray() (*RESPValue, error) {
	lengthStr, err := d.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return &RESPValue{Type: ARRAY, Null: true}, nil
	}

	array := make([]RESPValue, length)
	for i := 0; i < length; i++ {
		val, err := d.Decode()
		if err != nil {
			return nil, err
		}
		array[i] = *val
	}

	return &RESPValue{Type: ARRAY, Array: array}, nil
}

func (d *Decoder) decodeNull() (*RESPValue, error) {
	_, err := d.readLine()
	if err != nil {
		return nil, err
	}
	return &RESPValue{Type: NULL, Null: true}, nil
}

func (d *Decoder) decodeBoolean() (*RESPValue, error) {
	str, err := d.readLine()
	if err != nil {
		return nil, err
	}

	if str == "t" {
		return &RESPValue{Type: BOOLEAN, Bool: true}, nil
	} else if str == "f" {
		return &RESPValue{Type: BOOLEAN, Bool: false}, nil
	}

	return nil, errors.New("invalid boolean value: " + str)
}

func (d *Decoder) decodeDouble() (*RESPValue, error) {
	str, err := d.readLine()
	if err != nil {
		return nil, err
	}

	switch str {
	case "inf":
		return &RESPValue{Type: DOUBLE, Float: math.Inf(1)}, nil
	case "-inf":
		return &RESPValue{Type: DOUBLE, Float: math.Inf(-1)}, nil
	case "nan":
		return &RESPValue{Type: DOUBLE, Float: math.NaN()}, nil
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, err
	}
	return &RESPValue{Type: DOUBLE, Float: val}, nil
}

func (d *Decoder) decodeBigNumber() (*RESPValue, error) {
	str, err := d.readLine()
	if err != nil {
		return nil, err
	}

	bigInt := new(big.Int)
	_, ok := bigInt.SetString(str, 10)
	if !ok {
		return nil, errors.New("invalid big number: " + str)
	}

	return &RESPValue{Type: BIG_NUMBER, BigInt: bigInt}, nil
}

func (d *Decoder) decodeBulkError() (*RESPValue, error) {
	lengthStr, err := d.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(d.reader, data)
	if err != nil {
		return nil, err
	}

	d.readLine()
	return &RESPValue{Type: BULK_ERROR, Error: string(data)}, nil
}

func (d *Decoder) decodeVerbatimString() (*RESPValue, error) {
	lengthStr, err := d.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(d.reader, data)
	if err != nil {
		return nil, err
	}

	d.readLine()

	content := string(data)
	if len(content) >= 4 && content[3] == respTypeToChar[INTEGER] {
		return &RESPValue{Type: VERBATIM, String: content[4:]}, nil
	}

	return &RESPValue{Type: VERBATIM, String: content}, nil
}

func (d *Decoder) decodeMap() (*RESPValue, error) {
	lengthStr, err := d.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	m := make(map[string]RESPValue)
	for i := 0; i < length; i++ {
		key, err := d.Decode()
		if err != nil {
			return nil, err
		}

		value, err := d.Decode()
		if err != nil {
			return nil, err
		}

		m[key.String] = *value
	}

	return &RESPValue{Type: MAP, Map: m}, nil
}

func (d *Decoder) decodeAttribute() (*RESPValue, error) {
	lengthStr, err := d.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	m := make(map[string]RESPValue)
	for i := 0; i < length; i++ {
		key, err := d.Decode()
		if err != nil {
			return nil, err
		}

		value, err := d.Decode()
		if err != nil {
			return nil, err
		}

		m[key.String] = *value
	}

	return &RESPValue{Type: ATTRIBUTES, Map: m}, nil
}

func (d *Decoder) decodeSet() (*RESPValue, error) {
	lengthStr, err := d.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	array := make([]RESPValue, length)
	for i := 0; i < length; i++ {
		val, err := d.Decode()
		if err != nil {
			return nil, err
		}
		array[i] = *val
	}

	return &RESPValue{Type: SET, Array: array}, nil
}

func (d *Decoder) decodePush() (*RESPValue, error) {
	lengthStr, err := d.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}

	array := make([]RESPValue, length)
	for i := 0; i < length; i++ {
		val, err := d.Decode()
		if err != nil {
			return nil, err
		}
		array[i] = *val
	}

	return &RESPValue{Type: PUSH, Array: array}, nil
}

func DecodeString(input string) (*RESPValue, error) {
	decoder := NewDecoder(strings.NewReader(input))
	return decoder.Decode()
}
