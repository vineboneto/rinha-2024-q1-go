package services

import (
	"sync"

	"github.com/vineboneto/rinha-2024-q1-go/pg"
)

var (
	cliente404 = sync.Map{}
)

func FindCliente(id int64) bool {
	db := pg.GetDB()

	v, _ := cliente404.Load(id)

	if v == "404" {
		return false
	}

	if v == nil {
		type Out struct {
			ID int64 `gorm:"column:id"`
		}
		out := Out{}
		err := db.Raw("select id from clientes where id = ?", id).Scan(&out).Error

		if err != nil {
			cliente404.Store(id, "404")
			return false
		}

		if out.ID == 0 {
			cliente404.Store(id, "404")
			return false
		}

		cliente404.Store(id, "200")
	}

	return true
}
