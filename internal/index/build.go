package index

import (
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	// storage "github.com/sebzz2k2/vaultic/internal/storage"
	"github.com/sebzz2k2/vaultic/internal/wal"
	"github.com/sebzz2k2/vaultic/pkg/config"
)

type Index struct {
	filename string
	index    *sync.Map
	wal      *wal.WAL
}

func NewIndex(filename string, wal *wal.WAL) *Index {
	return &Index{filename: filename, wal: wal, index: &sync.Map{}}
}

func (idx *Index) bytesDecode(val []byte, decodedData *[]interface{}) {
	if len(val) < 4 {
		return
	}
	length := int(val[0])<<24 | int(val[1])<<16 | int(val[2])<<8 | int(val[3])

	if len(val) < length {
		return // guard against invalid/malformed data
	}

	entry := make([]byte, length)
	copy(entry, val[:length])

	decode, err := idx.wal.DecodeWAL(entry)
	if err != nil {
		log.Error().Err(err).Msg(config.ErrorDecodeData)
		return
	}
	*decodedData = append(*decodedData, decode)
	idx.bytesDecode(val[length:], decodedData)
}

func (idx *Index) BuildIndexes() error {
	fileBytes, err := os.ReadFile(idx.filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Msgf(config.ErrorNoFileFound, idx.filename)
			return nil
		}
		fmt.Println(err)
		return err
	}

	var decodedData []interface{}
	idx.bytesDecode(fileBytes, &decodedData)

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
			if idx.Exists(key) {
				idx.Del(key)
			}
			continue
		}
		idx.Set(key, start, offset)
	}

	return nil
}
