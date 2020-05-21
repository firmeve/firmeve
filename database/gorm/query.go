package gorm

import (
	"errors"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/reflect"
	"github.com/jinzhu/gorm"
	reflect2 "reflect"
)

type (
	gormQueryMagic struct {
		magic interface{}
		//Dto   interface{}
	}
)

func NewGORMQueryMagic(magic interface{}) (contract.GORMQueryMagic, error) {
	if reflect.KindElemType(reflect2.TypeOf(magic)) != reflect2.Struct {
		return nil, errors.New("magic not support the type")
	}
	//
	//if reflect.KindElemType(dto) != reflect2.Struct {
	//	panic(errors.New("dto not support the type"))
	//}

	return &gormQueryMagic{magic: magic}, nil
}

func NewGORMQueryMagicMust(magic interface{}) contract.GORMQueryMagic {
	m, err := NewGORMQueryMagic(magic)
	if err != nil {
		panic(err)
	}
	return m
}

func (g *gormQueryMagic) Query(db *gorm.DB) *gorm.DB {
	//dtoType := reflect2.TypeOf(g.Dto)
	//
	//dtoValue := reflect2.ValueOf(dto)
	selfType := reflect2.TypeOf(g.magic)
	selfValue := reflect2.ValueOf(g.magic)
	reflect.CallFieldType(selfType, func(i int, field reflect2.StructField) interface{} {
		fieldValue := reflect.CallOriginalFieldValue(selfValue, field.Name)
		if !fieldValue.IsZero() {
			methodName := `By` + field.Name
			if reflect.MethodExists(selfType, methodName) {
				db = reflect.CallMethodValue(selfValue, methodName, db, fieldValue)[0].(*gorm.DB)
			}
		}

		return nil
	})

	return db
}
