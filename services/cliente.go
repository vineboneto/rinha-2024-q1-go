package services

import (
	"github.com/vineboneto/rinha-2024-q1-go/pg"
)

func FindCliente(id int64) bool {
	db := pg.GetDB()

	ID := int64(0)
	err := db.Raw("select id from clientes where id = ?", id).Scan(&ID).Error

	if err != nil {
		return false
	}

	if ID == 0 {
		return false
	}

	return true
}
