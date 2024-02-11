package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/nagahshi/pos_go_desafio1/server/db/sqlite"
)

const HOST_API = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

// main - api para consulta de cotacao atual
func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

// cotacaoHandler
func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	// Chama a API para obter a cotação
	req, err := http.NewRequestWithContext(ctx, "GET", HOST_API, nil)
	if err != nil {
		if ctx.Err() != nil { // ctx nil ...
			log.Printf("[NewRequestWithContext] context error timout on request")
			http.Error(w, "Demorou muito para responder", http.StatusRequestTimeout)
			return
		}

		http.Error(w, "não foi possível buscar conversão do dolar", http.StatusUnprocessableEntity)
		return
	}

	// tenta realizar o request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			log.Printf("[request.do] context error timout on request")
			http.Error(w, "Demorou muito para responder", http.StatusRequestTimeout)
			return
		}

		http.Error(w, "não foi possível buscar conversão do dolar", http.StatusUnprocessableEntity)
		return
	}
	defer resp.Body.Close()

	// valida o retorno via status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "não foi possível buscar conversão do dolar", http.StatusUnprocessableEntity)
		return
	}

	// estrutura de resposta padrão da API
	var result struct {
		USDBRL struct {
			Bid string `json:"bid"`
		}
	}
	// decodifica o retorno usando estrutura padrão
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		http.Error(w, "não foi possível buscar conversão do dolar", http.StatusUnprocessableEntity)
		return
	}

	// // registra log
	err = newLogCotacao(result.USDBRL.Bid, &ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Envia a cotação para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // faz o assert mesmo que desnecessário no retorno
	json.NewEncoder(w).Encode(result)
}

// newLogCotacao - Registra a cotação no banco de dados SQLite
func newLogCotacao(bid string, ctx *context.Context) error {
	dbCtx, dbCancel := context.WithTimeout(*ctx, 10*time.Millisecond)
	defer dbCancel()

	// abre uma nova conexão SQLITE
	db := sqlite.NewConnDB("./cotacao.db")
	defer db.Close()

	_, err := db.ExecContext(dbCtx, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	if err != nil {
		if dbCtx.Err() != nil {
			log.Printf("[newLogCotacao] context error timout on register")
			return errors.New("tempo excedido")
		}

		return errors.New("não foi possível registrar conversão do dolar")
	}

	return nil
}
