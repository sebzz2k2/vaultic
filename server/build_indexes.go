package server

import (
	"os"

	"github.com/sebzz2k2/vaultic/kv_store/storage"
	"github.com/sebzz2k2/vaultic/logger"
	"github.com/sebzz2k2/vaultic/utils"
)

type IndexBuilder struct {
	filename  string
	delimiter byte
}

func NewIndexBuilder(filename string, delimiter byte) *IndexBuilder {
	return &IndexBuilder{filename: filename, delimiter: delimiter}
}

func bytesDecode(val []byte, decodedData *[]interface{}) {
	if len(val) < 4 {
		return
	}
	len := int(val[0])<<24 | int(val[1])<<16 | int(val[2])<<8 | int(val[3])
	entry := make([]byte, len)
	copy(entry, val[:len])
	decode, err := storage.DecodeData(entry)
	if err != nil {
		logger.Errorf("Error decoding data: %s", err.Error())
		return
	}
	*decodedData = append(*decodedData, decode)
	bytesDecode(val[len:], decodedData)
}

func (ib *IndexBuilder) BuildIndexes() error {
	file, err := os.Open(ib.filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fileBytes, err := os.ReadFile(ib.filename)
	if err != nil {
		return err
	}
	var decodedData []interface{}

	bytesDecode(fileBytes, &decodedData)

	offset := uint32(0)
	for _, val := range decodedData {
		if val == nil {
			continue
		}
		offset += val.(map[string]interface{})["totalLength"].(uint32)
		valLen := val.(map[string]interface{})["valueLen"].(uint32)
		start := offset - valLen
		key := val.(map[string]interface{})["key"].(string)
		utils.SetIndexKey(key, start, offset)
	}
	utils.PrintIndexMap()
	return nil
}
