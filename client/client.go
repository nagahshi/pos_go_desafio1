package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const serverURL = "http://localhost:8080/cotacao"

// main - efetua consulta na api de cotacoes
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		if ctx.Err() != nil {
			log.Printf("[NewRequestWithContext] context error timout on request")
		}

		fmt.Println("Erro ao preparar request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			log.Printf("[request.do] context error timout on request")
		}

		fmt.Println("Erro ao realizar request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Não foi possível ler o corpo da resposta:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("api respondeu de forma inesperada:", resp.StatusCode)
		return
	}

	var request struct {
		USDBRL struct {
			Bid string `json:"bid"`
		}
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("Não foi possível ler o retorno da api:", err)
		return
	}

	// registrando cotacao
	err = registroLogCotacao([]byte(fmt.Sprintf("Dólar: %s\n", request.USDBRL.Bid)))
	if err != nil {
		fmt.Println("Não foi possível registrar o log:", err)
		return
	}
}

// registroLogCotacao - registro de cotacoes
func registroLogCotacao(data []byte) (err error) {
	// cria e atribui permissões no arquivo
	f, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// faz o append de informações
	if _, err := f.Write(data); err != nil {
		return err
	}

	// fecha o arquivo
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
