FROM golang:1.25-alpine

WORKDIR /app

# Copia os arquivos do projeto
COPY go.mod ./
COPY server.go ./

# Compila a aplicação
RUN go build -o server server.go

# Expõe a porta 8080
EXPOSE 8080

# Executa o servidor
CMD ["./server"]