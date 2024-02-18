package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/vineboneto/rinha-2024-q1-go/pg"
	"github.com/vineboneto/rinha-2024-q1-go/services"
)

func main() {
	r := chi.NewRouter()

	defer pg.CloseDB()

	pg.GetDB()

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Post("/clientes/{id}/transacoes", services.TransacaoController)
	r.Get("/clientes/{id}/extrato", services.Extrato)

	log.Printf("Servidor iniciado na porta %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)

	if err != nil {
		log.Printf("Erro ao iniciar o servidor: %s\n", err)
	}
}
