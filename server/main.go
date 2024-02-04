package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	cotacao, err := getCotacao(&ctxAPI)
	if err != nil {
		http.Error(w, "Erro ao obter a cotação", http.StatusInternalServerError)
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
		return "", fmt.Errorf("Erro ao criar requisição para a API: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if apiCtx.Err() != nil {
			fmt.Println("err de contexto request")
		}
		return "", fmt.Errorf("Erro ao chamar a API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status não OK: %s", resp.Status)
	}

	var result map[string]map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("Erro ao decodificar resposta da API: %v", err)
	}

	// registra log

	return "", nil
}
