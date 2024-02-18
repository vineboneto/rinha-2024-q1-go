package services

import (
	"github.com/vineboneto/rinha-2024-q1-go/pg"
)

// var (
// 	cliente404 sync.Map
// )

func FindCliente(id int64) bool {
	db := pg.GetDB()

	// v, _ := cliente404.Load(id)
	ID := int64(0)
	err := db.Raw("select id from clientes where id = ?", id).Scan(&ID).Error

	if err != nil {
		// cliente404.Store(id, "404")
		return false
	}

	if ID == 0 {
		// cliente404.Store(id, "404")
		return false
	}

	// if v == "404" {
	// 	return
	// }

	// if v == nil {

	// 	cliente404.Store(id, "200")
	// }

	return true
}
