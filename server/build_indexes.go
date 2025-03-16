package server

import (
	"bufio"
	"os"

	"github.com/sebzz2k2/vaultic/utils"
)

type IndexBuilder struct {
	filename  string
	delimiter byte
}

func NewIndexBuilder(filename string, delimiter byte) *IndexBuilder {
	return &IndexBuilder{filename: filename, delimiter: delimiter}
}

func (ib *IndexBuilder) BuildIndexes() error {
	file, err := os.Open(ib.filename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var key []byte
	i, ignore := 0, false
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		i++
		if ignore {
			if b == '\n' {
				ignore = false
			}
			continue
		}
		if b == ib.delimiter {
			utils.SetIndexKey(string(key), i+1)
			ignore, key = true, nil
			continue
		}
		key = append(key, b)
	}
	return nil
}
