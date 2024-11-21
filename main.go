package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	cep := os.Args[1]
	url1 := "https://brasilapi.com.br/api/cep/v1/" + cep
	url2 := "http://viacep.com.br/ws/" + cep + "/json/"
	chBr := make(chan BrasilAPIResponse)
	chVia := make(chan ViaCepResponse)

	go func() {
		req, err := http.Get(url2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		defer req.Body.Close()
		var response ViaCepResponse
		err = json.NewDecoder(req.Body).Decode(&response)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		chVia <- response
	}()

	go func() {
		req, err := http.Get(url1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		defer req.Body.Close()
		var response BrasilAPIResponse
		err = json.NewDecoder(req.Body).Decode(&response)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		chBr <- response
	}()

	select {
	case br := <-chBr:
		fmt.Printf("Brasil-API: %v\n", br)
	case via := <-chVia:
		fmt.Printf("ViaCEP-API: %v\n", via)
	case <-time.After(1 * time.Second):
		log.Fatal("Timeout")
	}
}
