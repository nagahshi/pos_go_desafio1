package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctxAPI, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
	defer cancel()

	// Chama a função para obter a cotação
	err := getCotacao(&ctxAPI)
	if err != nil {
		http.Error(w, "Erro ao obter a cotação", http.StatusInternalServerError)
		return
	}

	// Envia a cotação para o cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("")
}

func getCotacao(ctx *context.Context) error {
	return nil
}
