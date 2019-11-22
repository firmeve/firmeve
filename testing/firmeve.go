package testing

import "github.com/firmeve/firmeve"

func TestingModeFirmeve() *firmeve.Firmeve  {
	return firmeve.New(firmeve.WithMode(firmeve.ModeTesting))
}