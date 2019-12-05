package page

import (
	"github.com/firmeve/firmeve/converter/resource"
	"github.com/jinzhu/gorm"
	"math"
)

type GORM struct {
	db *gorm.DB
	//option *resource.PageOption
	page  uint
	limit uint
	//done  chan error
	count uint
	info  *resource.PaginatorInfo
}

func NewGorm(db *gorm.DB, page, limit uint) *GORM {
	return &GORM{
		db:    db,
		page:  page,
		limit: limit,
		count: 0,
		//done:  make(chan error, 0),
		//info: &resource.PaginatorInfo{},
	}
}

func (g *GORM) Info() *resource.PaginatorInfo {
	return g.info
}

func (g *GORM) skip() uint {
	if g.page == 1 {
		return 0
	} else {
		return (g.page - 1) * g.limit
	}
}

func (g *GORM) Resource(v interface{}) (interface{}, error) {
	done := make(chan error, 1)

	go g.Count(v, done)

	g.db.Limit(g.limit).Offset(g.skip()).Find(v)
	err := <-done
	if err != nil {
		return nil, err
	}

	g.info = &resource.PaginatorInfo{
		Total:     g.count,
		TotalPage: uint(math.Ceil(float64(g.count) / float64(g.limit))),
		Limit:     g.limit,
		Offset:    g.skip(),
		Next:      g.page + 1,
		Prev:      g.page - 1,
	}

	return v, nil
}

func (g *GORM) Count(v interface{}, done chan<- error) {
	if err := g.db.Model(v).Count(&g.count).Error; err != nil {
		done <- err
	} else {
		done <- nil
	}
}
