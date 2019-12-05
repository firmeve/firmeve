package page

import (
	"fmt"
	bootstrap2 "github.com/firmeve/firmeve/bootstrap"
	"github.com/firmeve/firmeve/support/path"
	"github.com/firmeve/firmeve/support/strings"
	"github.com/firmeve/firmeve/testing/bootstrap"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Test struct {
	Name string
	Uuid string
}

var (
	//f  *firmeve.Firmeve
	db *gorm.DB
)

func init() {
	//f = testing2.TestingModeFirmeve()
	//f.Register(f.Make(new(database.Provider)).(*database.Provider))
	b := bootstrap.Bootstrap(bootstrap2.WithConfigPath(path.RunRelative("../../../testdata/config")))
	db = b.Firmeve.Get(`db.connection`).(*gorm.DB)
}

func TestNewGorm(t *testing.T) {
	db.Exec("TRUNCATE table `tests`")
	for i := 0; i <= 35; i++ {
		db.Table(`tests`).Create(&Test{
			Name: strings.Rand(30),
			Uuid: strings.Rand(36),
		})
	}

	gorm := NewGorm(db,1,20)
	var tests []*Test
	v,err := gorm.Resource(&tests)
	assert.Nil(t,err)
	fmt.Println(tests)
	fmt.Println(v)
}
