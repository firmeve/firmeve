package resource

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/database"
	"github.com/firmeve/firmeve/support/path"
	"github.com/firmeve/firmeve/support/strings"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/jinzhu/gorm"
	"github.com/ulule/paging"
	"net/http"
	"testing"
	"fmt"
)

type (
	Test struct {
		Name string
		Uuid string
	}
)

func (t Test) TableName() string  {
	return `tests`
}

var (
	db *gorm.DB
)

func init() {
	f := testing2.TestingModeFirmeve()
	f.Bind(`config`, config.New(path.RunRelative("../../testdata/config")), container.WithShare(true))
	f.Register(f.Make(new(database.Provider)).(firmeve.Provider))
	f.Boot()
	db = f.Get(`db.connection`).(*gorm.DB)
}

func TestNewPaginator(t *testing.T) {
	db.Exec("TRUNCATE table `tests`")
	for i := 0; i <= 35; i++ {
		db.Table(`tests`).Create(&Test{
			Name: strings.Rand(30),
			Uuid: strings.Rand(36),
		})
	}

	var tests []*Test
	store,_ := paging.NewGORMStore(db,&tests)
	paginator := NewPaginator(store, testing2.NewMockRequest(http.MethodGet,"http://firmeve.com/?limit=5&offset=1","").Request,&Option{
		Transformer:nil,
		Fields:[]string{"name","uuid"},
	},nil)

	fmt.Println(paginator.CollectionData())
	fmt.Println(paginator.Meta())
	fmt.Println(paginator.Link())
}
