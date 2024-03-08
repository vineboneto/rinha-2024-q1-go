package services

import (
	"github.com/vineboneto/rinha-2024-q1-go/pg"
)

func FindCliente(id int64) bool {
	db := pg.GetDB()

	out := int64(0)
	err := db.Raw("select id from clientes where id = ?", id).Scan(&out).Error

	if err != nil {
		return false
	}

	if out == 0 {
		return false
	}

	return true
}
