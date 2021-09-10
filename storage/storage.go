package storage

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/storage/disk"
)

const Default = disk.QiNiuDisk

type (
	Storage struct {
		config *Configuration
		disks  map[string]contract.Storage
	}

	Configuration struct {
		QiNiu *disk.QiNiuConfiguration `json:"qiniu" yaml:"qiniu"`
	}
)

func New(config *Configuration) *Storage {
	return &Storage{config: config, disks: make(map[string]contract.Storage, 1)}
}

func (s Storage) RegisterDisk(name string, storage contract.Storage) {
	s.disks[name] = storage
}

func (s Storage) Disk(name string) contract.Storage {
	if v, ok := s.disks[name]; ok {
		return v
	}
	panic(`disk [` + name + `] not found`)
}
