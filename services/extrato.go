package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/vineboneto/rinha-2024-q1-go/pg"
)

type ExtratoOutput struct {
	Saldo   json.RawMessage `gorm:"column:saldo" json:"saldo"`
	Extrato json.RawMessage `gorm:"column:extrato" json:"ultimas_transacoes"`
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
				json_build_object(
						'total', c.saldo,
						'limite', c.limite,
						'data_extrato', TO_CHAR(NOW() AT TIME ZONE 'utc', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
				) as saldo,
				COALESCE(
					(
							SELECT JSON_AGG(f.*)
							FROM (
									SELECT 
											t.valor,
											t.descricao,
											t.tipo,
											TO_CHAR(t.realizada_em AT TIME ZONE 'utc', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS realizada_em
									FROM transacoes t
									WHERE t.id_cliente = c.id
									ORDER BY t.realizada_em DESC
									LIMIT 10
							) AS f
					),
					'[]'
			) AS extrato
		from clientes c where c.id = ?
	`

	output := ExtratoOutput{}

	db.Raw(sql, clienteId).Scan(&output)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
