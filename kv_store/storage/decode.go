package storage

import (
	"encoding/binary"
	"errors"

	"github.com/sebzz2k2/vaultic/utils"
)

func DecodeFlags(encoded byte) map[string]interface{} {
	flags := make([]bool, 8)
	for i := 0; i < 8; i++ {
		flags[i] = (encoded & (1 << i)) != 0
	}
	return map[string]interface{}{
		"deleted":    flags[0],
		"compressed": flags[1],
		"reserved":   flags[3:7],
	}
}
func DecodeData(encoded []byte) (map[string]interface{}, error) {
	if len(encoded) < 4 {
		return nil, errors.New("invalid data: insufficient length")
	}

	totalLength := binary.BigEndian.Uint32(encoded[:4])
	if len(encoded) != int(totalLength) {
		return nil, errors.New("data length mismatch")
	}

	version := encoded[4]
	flags := encoded[5]
	decodedFlags := DecodeFlags(flags)
	keyValCRC := binary.BigEndian.Uint32(encoded[6:10])
	ts := binary.BigEndian.Uint64(encoded[10:18])

	keyLen := binary.BigEndian.Uint16(encoded[18:20])
	valueLen := binary.BigEndian.Uint32(encoded[20:24])

	if int(24+uint32(keyLen)+valueLen) != len(encoded) {
		return nil, errors.New("mismatched key/value lengths")
	}

	key := string(encoded[24 : 24+keyLen])
	value := string(encoded[24+uint32(keyLen) : 24+uint32(keyLen)+valueLen])

	// Verify CRC
	computedCRC := utils.Crc32(key + value)
	if computedCRC != keyValCRC {
		return nil, errors.New("CRC check failed")
	}
	return map[string]interface{}{
		"totalLength": totalLength,
		"version":     version,
		"flags":       decodedFlags,
		"key":         key,
		"keyLen":      keyLen,
		"valueLen":    valueLen,
		"val":         value,
		"ts":          ts,
	}, nil
}
