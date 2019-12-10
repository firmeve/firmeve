package resource

import (
	"fmt"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/database"
	"github.com/firmeve/firmeve/support/path"
	"github.com/firmeve/firmeve/support/strings"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/paging"
	"net/http"
	"testing"
)

type (
	Test struct {
		Id   uint
		Name string
		Uuid string
	}
)

func (t Test) TableName() string {
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

func TestNewPaginator2(t *testing.T) {
	db.Exec("TRUNCATE table `tests`")
	for i := 0; i <= 35; i++ {
		db.Table(`tests`).Create(&Test{
			Name: strings.Rand(30),
			Uuid: strings.Rand(36),
		})
	}

	db = db.New().Order("id asc").Debug()

	var tests []*Test
	limit := int64(15) //strconv.ParseInt(req.URL.Query().Get("limit"), 10, 64)
	pageOption := &paging.Options{
		DefaultLimit:  limit,
		MaxLimit:      limit + 10,
		LimitKeyName:  "limit",
		OffsetKeyName: "offset",
	}
	option := &Option{
		Transformer: nil,
		Fields:      []string{"id", "name", "uuid"},
	}
	store, _ := paging.NewGORMStore(db, &tests)

	req := testing2.NewMockRequest(http.MethodGet, "http://firmeve.com/any/testing?any_param=1&limit=15000&offset=1000", "").Request
	paginator := NewPaginator(store, option, req, pageOption)
	fmt.Println("==================")
	fmt.Printf("%#v\n",paginator.CollectionData())
	meta := paginator.Meta()
	link := paginator.Link()
	fmt.Printf("%#v\n",meta)
	fmt.Println("==================")
	fmt.Printf("%#v\n",link)

	//assert.Equal(t, int64(1), meta[`current_page`].(int64))
	//assert.Equal(t, int64(36), meta[`total`].(int64))
	//assert.Equal(t, int64(15), meta[`per_page`].(int64))
	//assert.Equal(t, int64(1), meta[`from`].(int64))
	//assert.Equal(t, int64(15), meta[`to`].(int64))
	//assert.Equal(t, int64(3), meta[`last_page`].(int64))
	//assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=0&any_param=1", link["first"])
	//assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=30&any_param=1", link["last"])
	//assert.Equal(t, "", link["prev"])
	//assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=15&any_param=1", link["next"])
}

func TestNewPaginator(t *testing.T) {
	db.Exec("TRUNCATE table `tests`")
	for i := 0; i <= 35; i++ {
		db.Table(`tests`).Create(&Test{
			Name: strings.Rand(30),
			Uuid: strings.Rand(36),
		})
	}

	db = db.New().Order("id asc").Debug()

	//
	var tests []*Test
	limit := int64(15) //strconv.ParseInt(req.URL.Query().Get("limit"), 10, 64)
	pageOption := &paging.Options{
		DefaultLimit:  limit,
		MaxLimit:      limit + 10,
		LimitKeyName:  "limit",
		OffsetKeyName: "offset",
	}
	option := &Option{
		Transformer: nil,
		Fields:      []string{"id", "name", "uuid"},
	}
	store, _ := paging.NewGORMStore(db, &tests)

	req := testing2.NewMockRequest(http.MethodGet, "http://firmeve.com/any/testing?any_param=1&limit=15&offset=0", "").Request
	paginator := NewPaginator(store, option, req, pageOption)
	assert.Equal(t, 15, len(paginator.CollectionData()))
	assert.Equal(t, uint(1), paginator.CollectionData()[0][`id`].(uint))
	meta := paginator.Meta()
	link := paginator.Link()
	assert.Equal(t, int64(1), meta[`current_page`].(int64))
	assert.Equal(t, int64(36), meta[`total`].(int64))
	assert.Equal(t, int64(15), meta[`per_page`].(int64))
	assert.Equal(t, int64(1), meta[`from`].(int64))
	assert.Equal(t, int64(15), meta[`to`].(int64))
	assert.Equal(t, int64(3), meta[`last_page`].(int64))
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=0&any_param=1", link["first"])
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=30&any_param=1", link["last"])
	assert.Equal(t, "", link["prev"])
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=15&any_param=1", link["next"])
	//
	//
	//store1, _ := paging.NewGORMStore(db.Debug(), &tests)
	req1 := testing2.NewMockRequest(http.MethodGet, "http://firmeve.com/any/testing?any_param=1&limit=15&offset=15", "").Request
	paginator1 := NewPaginator(store, option, req1, pageOption)
	assert.Equal(t, 15, len(paginator1.CollectionData()))
	assert.Equal(t, uint(16), paginator1.CollectionData()[0][`id`].(uint))
	meta1 := paginator1.Meta()
	link1 := paginator1.Link()
	assert.Equal(t, int64(2), meta1[`current_page`].(int64))
	assert.Equal(t, int64(36), meta1[`total`].(int64))
	assert.Equal(t, int64(15), meta1[`per_page`].(int64))
	assert.Equal(t, int64(16), meta1[`from`].(int64))
	assert.Equal(t, int64(30), meta1[`to`].(int64))
	assert.Equal(t, int64(3), meta1[`last_page`].(int64))
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=0&any_param=1", link1["first"])
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=30&any_param=1", link1["last"])
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=0&any_param=1", link1["prev"])
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=30&any_param=1", link1["next"])
	//
	req2 := testing2.NewMockRequest(http.MethodGet, "http://firmeve.com/any/testing?any_param=1&limit=15&offset=30", "").Request
	paginator2 := NewPaginator(store, option, req2, pageOption)
	assert.Equal(t, 6, len(paginator2.CollectionData()))
	assert.Equal(t, uint(31), paginator2.CollectionData()[0][`id`].(uint))
	meta2 := paginator2.Meta()
	link2 := paginator2.Link()
	assert.Equal(t, int64(3), meta2[`current_page`].(int64))
	assert.Equal(t, int64(36), meta2[`total`].(int64))
	assert.Equal(t, int64(15), meta2[`per_page`].(int64))
	assert.Equal(t, int64(31), meta2[`from`].(int64))
	assert.Equal(t, int64(36), meta2[`to`].(int64))
	assert.Equal(t, int64(3), meta2[`last_page`].(int64))
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=0&any_param=1", link2["first"])
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=30&any_param=1", link2["last"])
	assert.Equal(t, "http://firmeve.com/any/testing?limit=15&offset=15&any_param=1", link2["prev"])
	assert.Equal(t, "", link2["next"])



	//req2 := testing2.NewMockRequest(http.MethodGet, "http://firmeve.com/any/testing?any_param=1&limit=15&offset=30", "").Request
	//paginator2 := NewPaginator(store, option, req2, pageOption)
	//assert.Equal(t, 5, len(paginator2.CollectionData()))
	//assert.Equal(t, uint(31), paginator2.CollectionData()[0][`id`].(uint))
	//
	//fmt.Printf("%#v\n", paginator2.CollectionData())
	//fmt.Println(paginator.Meta())
	//fmt.Println(paginator.Link())
}
