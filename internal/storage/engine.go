package storage

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/sebzz2k2/vaultic/internal/index"
	"github.com/sebzz2k2/vaultic/internal/protocol"
	"github.com/sebzz2k2/vaultic/internal/wal"
)

type StorageEngine struct {
	idx      *index.Index
	Protocol *protocol.Protocol
	wal      *wal.WAL
}

func NewStorageEngine() (*StorageEngine, error) {
	wal := wal.NewWAL()
	idx := index.NewIndex("vaultic", wal)

	log.Info().Msg("Building indexes")
	if err := idx.BuildIndexes(); err != nil {
		return nil, fmt.Errorf("failed to build indexes: %w", err)
	}
	log.Info().Msg("Indexes built successfully")

	return &StorageEngine{
		wal:      wal,
		idx:      idx,
		Protocol: protocol.NewProtocol(wal),
	}, nil
}

func (se *StorageEngine) Get()           {}
func (se *StorageEngine) Set()           {}
func (se *StorageEngine) Delete()        {}
func (se *StorageEngine) Exists() bool   { return false }
func (se *StorageEngine) Keys() []string { return []string{} }
func (se *StorageEngine) Close() error   { return nil }
