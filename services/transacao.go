package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/vineboneto/rinha-2024-q1-go/pg"
)

type TransacaoInput struct {
	Valor     int64  `json:"valor"`
	Descricao string `json:"descricao"`
	Tipo      string `json:"tipo"`
}

type UpdateOutput struct {
	Saldo  int64 `gorm:"column:saldo" json:"saldo"`
	Limite int64 `gorm:"column:limite" json:"limite"`
}

func TransacaoController(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	clienteId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var input TransacaoInput

	err = json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	isValid := (input.Tipo == "c" || input.Tipo == "d") && input.Valor > 0 && len(input.Descricao) > 0 && len(input.Descricao) <= 10

	if !isValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	db := pg.GetDB()

	exist := FindCliente(clienteId)

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	valorIncrementado := input.Valor

	if input.Tipo == "d" {
		valorIncrementado = -valorIncrementado
	}

	output := UpdateOutput{}

	result := db.Raw(`
		update clientes 
		set saldo = saldo + ? where id = ? and (saldo + ?) * -1 <= limite 
		returning saldo, limite
	`, valorIncrementado, clienteId, valorIncrementado).
		Scan(&output)

	if result.Error != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	db.Exec("insert into transacoes (id_cliente, valor, descricao, tipo) values (?, ?, ?, ?)", clienteId, input.Valor, input.Descricao, input.Tipo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
