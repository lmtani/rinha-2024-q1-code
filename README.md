# rinha-de-backend-2024q1


## Development Mode

There is a resource limitation as part of the tests (see [Rinha de Backend - 2024/Q1 / Restrições de CPU/Memória](https://github.com/zanfranceschi/rinha-de-backend-2024-q1?tab=readme-ov-file#restri%C3%A7%C3%B5es-de-cpumem%C3%B3ria)).

So, in development mode, we're overwriting the limits to be able to have hot-reload on:
```sh
docker-compose -f docker-compose.yml -f docker-compose-dev.yml up db api01
```
