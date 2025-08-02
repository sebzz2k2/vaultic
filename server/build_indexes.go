package server

import (
	"os"

	"github.com/rs/zerolog/log"
	storage "github.com/sebzz2k2/vaultic/kv_store"
	"github.com/sebzz2k2/vaultic/pkg/config"
	"github.com/sebzz2k2/vaultic/utils"
)

type IndexBuilder struct {
	filename string
}

func NewIndexBuilder(filename string) *IndexBuilder {
	return &IndexBuilder{filename: filename}
}

func bytesDecode(val []byte, decodedData *[]interface{}) {
	if len(val) < 4 {
		return
	}
	length := int(val[0])<<24 | int(val[1])<<16 | int(val[2])<<8 | int(val[3])

	if len(val) < length {
		return // guard against invalid/malformed data
	}

	entry := make([]byte, length)
	copy(entry, val[:length])

	decode, err := storage.DecodeWAL(entry)
	if err != nil {
		log.Error().Err(err).Msg(config.ErrorDecodeData)
		return
	}
	*decodedData = append(*decodedData, decode)
	bytesDecode(val[length:], decodedData)
}

func (ib *IndexBuilder) BuildIndexes() error {
	fileBytes, err := os.ReadFile(ib.filename)
	if err != nil {
		return err
	}

	var decodedData []interface{}
	bytesDecode(fileBytes, &decodedData)

	offset := uint32(0)
	for _, v := range decodedData {
		entry, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		key := entry["key"].(string)
		valLen := entry["valueLen"].(uint32)
		totalLen := entry["totalLength"].(uint32)
		offset += totalLen
		start := offset - valLen

		flags := entry["flags"].(map[string]interface{})
		if flags["deleted"].(bool) {
			if utils.IsPresent(key) {
				utils.DeleteIndexKey(key)
			}
			continue
		}

		utils.SetIndexKey(key, start, offset)
	}

	return nil
}
