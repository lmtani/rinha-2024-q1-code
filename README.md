# rinha-de-backend-2024q1


```sh
# Desenvolvimento local
## Suba o banco de dados
docker-compose up -d db

## Suba o servidor
cd rinha
DB_HOSTNAME=postgres://admin:123@localhost/rinha go run .
```