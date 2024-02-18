# rinha-de-backend-2024q1


```sh
# Desenvolvimento local
## Suba o banco de dados
docker-compose up -d db

## Suba o servidor
go build -o main ./cmd/server
DB_HOSTNAME=postgres://admin:123@localhost/rinha ./main
```