package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type retorno struct {
	URL  string
	body string
}

func callAPI(ctx context.Context, URL string, ch chan<- retorno) error {

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		log.Panicln(err)
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Erro na chamada da API %s: %v\n", URL, err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	ch <- retorno{
		URL:  URL,
		body: string(body),
	}
	return nil

}

func main() {

	// CEP pode ser passado como flag
	cep := flag.String("c", "80530000", "CEP para realizar a consulta")
	flag.Parse()

	// Slice de string das APIs a serem consultadas
	arr := []string{
		fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", *cep),
		fmt.Sprintf("http://viacep.com.br/ws/%s/json/", *cep),
	}

	// Channel de struct de retorno
	ch := make(chan retorno)

	// Context que será passado a todas as Go Routines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, URL := range arr {
		go func() {
			callAPI(ctx, URL, ch)
		}()
	}

	select {
	case r := <-ch:
		fmt.Printf("Retorno da API: %s\n%s\n", r.URL, r.body)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout - Nenhuma API retornou a tempo")
	}

}
