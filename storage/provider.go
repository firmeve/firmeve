package storage

import (
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/storage/disk"
)

type Provider struct {
	kernel.BaseProvider
}

func (p Provider) Name() string {
	return `storage`
}

func (p *Provider) Register() {
	var config = new(Configuration)
	p.Config.BindItem(`storage`, config)
	storage := New(config)
	storage.RegisterDisk(disk.QiNiuDisk, disk.NewQiNiu(config.QiNiu))
	p.Bind(`storage`, storage)
	p.Bind(`storage.qiniu`, storage.Disk(disk.QiNiuDisk))
}

func (p Provider) Boot() {
}
