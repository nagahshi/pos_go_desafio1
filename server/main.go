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

// main - api para consulta de cotacao atual
func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctxAPI, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
	defer cancel()

	// Chama a função para obter a cotação
	cotacao, err := getCotacao(&ctxAPI)
	if err != nil {
		if ctxAPI.Err() != nil {
			log.Printf("[cotacaoHandler] context error timout on handler")
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Envia a cotação para o cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao)
}

// getCotacao - busca cotacoes
func getCotacao(ctx *context.Context) (string, error) {
	apiCtx, apiCancel := context.WithTimeout(*ctx, 200*time.Millisecond)
	defer apiCancel()

	// Chama a API para obter a cotação
	req, err := http.NewRequestWithContext(apiCtx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", errors.New("não foi possível buscar conversão do dolar")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if apiCtx.Err() != nil {
			log.Printf("[getCotacao] context error timout on request")
			return "", errors.New("tempo excedido")
		}
		return "", errors.New("não foi possível buscar conversão do dolar")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("não foi possível buscar conversão do dolar")
	}

	var result map[string]map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", errors.New("não foi possível buscar conversão do dolar")
	}

	var bid string
	if result["USDBRL"] != nil && result["USDBRL"]["bid"] != "" {
		bid = result["USDBRL"]["bid"]
	}

	// // registra log
	err = newLogCotacao(bid, ctx)
	if err != nil {
		return "", err
	}

	return bid, nil
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
