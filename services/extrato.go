package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/vineboneto/rinha-2024-q1-go/pg"
)

type ExtratoOutput struct {
	Limite      int64           `gorm:"column:limite" json:"limite"`
	Saldo       int64           `gorm:"column:saldo" json:"saldo"`
	DataExtrato string          `json:"data_extrato" gorm:"-"`
	Extrato     json.RawMessage `gorm:"column:extrato" json:"ultimas_transacoes"`
}

func Extrato(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	clienteId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	db := pg.GetDB()

	exist := FindCliente(clienteId)

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sql := `
			select
			c.limite,
			c.saldo,
			(select json_agg(f.*) from (
				select 
						t.valor,
						t.descricao,
						t.tipo,
						to_char(t.realizada_em AT TIME ZONE 'utc', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as realizada_em
					from transacoes t
					where t.id_cliente = c.id
					order by t.realizada_em desc limit 10
				) as f
			) as extrato
		from clientes c where c.id = ?
	`

	output := ExtratoOutput{}

	db.Raw(sql, clienteId).Scan(&output)

	type Output struct {
		Saldo struct {
			Total  int64 `json:"total"`
			Limite int64 `json:"limite"`
		} `json:"saldo"`
		UltimasTransacoes json.RawMessage `json:"ultimas_transacoes"`
	}

	outJSON := Output{}

	outJSON.UltimasTransacoes = output.Extrato

	output.DataExtrato = time.Now().UTC().Format(ISO_8601)

	if outJSON.UltimasTransacoes == nil {
		outJSON.UltimasTransacoes = json.RawMessage("[]")
	}

	outJSON.Saldo.Total = output.Saldo
	outJSON.Saldo.Limite = output.Limite

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(outJSON)
}
