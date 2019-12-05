package resource

import (
	//"github.com/firmeve/firmeve/support/reflect"
	//"github.com/jinzhu/gorm"
	"github.com/ulule/paging"
	"net/http"
)

type Paginator struct {
	//db         *gorm.DB
	request    *http.Request
	resource   *paging.GORMStore
	option     *Option
	pageOption *paging.Options
	meta       Meta
	link       Link
}

func NewPaginator(resource *paging.GORMStore, request *http.Request, option *Option, pageOption *paging.Options) *Paginator {
	return &Paginator{
		//db:         db,
		request:    request,
		resource:   resource,//reflect.SliceInterface(reflect2.ValueOf(resource)),
		pageOption: pageOption,
		option:     option,
		meta:       make(Meta, 0),
		link:       make(Link, 0),
	}
}

func (p *Paginator) CollectionData() DataCollection {
	//store, err := paging.NewGORMStore(p.db, &p.resource)
	//if err != nil {
	//	panic(err)
	//}

	paginator, _ := paging.NewOffsetPaginator(p.resource, p.request, p.pageOption)

	err := paginator.Page()
	if err != nil {
		panic(err)
	}

	p.SetLink(Link{
		"prev":  paginator.PreviousURI.String,
		"next":  paginator.NextURI.String,
		"first": "",
		"last":  "",
	})

	p.SetMeta(Meta{
		"current_page": 1,               //当前页
		"total":        paginator.Count, //总条数
		"per_page":     "",              //每页多少条
		"from":         "",              //从多少条
		"last_page":    "",
		"to":           "", //到多少条
		"path":         "",
	})

	return NewCollection(p.resource.GetItems().([]interface{}), p.option).CollectionData()
}

func (p *Paginator) SetLink(links Link) {
	p.link = links
}

func (p *Paginator) Link() Link {
	return p.link
}

func (p *Paginator) SetMeta(meta Meta) {
	p.meta = meta
}

func (p *Paginator) Meta() Meta {
	return p.meta
}
