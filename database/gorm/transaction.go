package gorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func Transaction(db *gorm.DB, fn func(db *gorm.DB) (interface{}, error)) (result interface{}, tError error) {
	var (
		tx *gorm.DB
	)

	tx = db.Begin()

	defer func() {
		if tError := recover(); tError != nil {
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				tError = fmt.Errorf("transaction rollabck execute error %w", rollbackErr)
			}

			return
		}

		if commitErr := tx.Commit().Error; commitErr != nil {
			tError = fmt.Errorf("transaction commit execute error %w", commitErr)
		}
	}()

	result, tError = fn(tx)

	// 这里不能直接返回 defer的特性
	// 函数的返回过程是这样  先给返回赋值->调用defer->返回到调用函数中
	// 如果 return tError 则会先赋值 tError 此时 defer 还没有执行 也就是 nil
	// 参见 https://tiancaiamao.gitbooks.io/go-internals/content/zh/03.4.html
	return
}
