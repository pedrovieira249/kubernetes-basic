package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var startedAt time.Time = time.Now()

// Parara rodar o servidor, utilize o comando: go run server.go
func main() {
	http.HandleFunc("/", Hello)
	http.HandleFunc("/healthz", Healthz)
	http.HandleFunc("/config", ConfigMap)
	http.HandleFunc("/secret", Secret)
	http.ListenAndServe(":8080", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	saudacao := os.Getenv("SAUDACAO")
	mensagem := os.Getenv("MENSAGEM")
	w.Write([]byte("<h1>" + saudacao + "</h1>"))
	w.Write([]byte("<h2>" + mensagem + "</h2>"))
}

func ConfigMap(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("config/myconfig.txt")
	if err != nil {
		http.Error(w, "Erro ao ler o arquivo de configuração", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("<pre>" + string(data) + "</pre>"))
}

func Secret(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	w.Write([]byte("<h2>" + username + "</h2>"))
	w.Write([]byte("<h2>" + password + "</h2>"))
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	duration := time.Since(startedAt)
	if duration.Seconds() <= 10 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Tempo de atividade excedido: %.0f seconds", duration.Seconds())))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
