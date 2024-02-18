package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmtani/rinha-2024-q1-code/internal/services"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var dbpool *pgxpool.Pool

type Server struct {
	dbpool  *pgxpool.Pool
	service *services.Service
}

func NewServer(dbpool *pgxpool.Pool) *Server {
	return &Server{dbpool: dbpool, service: services.NewService(dbpool)}

}

func main() {
	ctx := context.Background()
	dbpool = initializeDatabase(ctx, dbpool)
	defer dbpool.Close()

	server := NewServer(dbpool)

	fmt.Println("Server running on port 8080")

	router := routing.New()
	router.Get("/clientes/<id>/extrato", server.StatementHandler)
	router.Post("/clientes/<id>/transacoes", server.TransactionsHandler)

	log.Fatal(fasthttp.ListenAndServe("0.0.0.0:8080", router.HandleRequest))
}

func (s *Server) StatementHandler(c *routing.Context) error {
	clientID, err := parseClientID(c.Param("id"))
	if err != nil {
		return respondWithError(c, "Invalid client ID", fasthttp.StatusNotFound)
	}

	r, err := s.service.HandleGetStatement(clientID)
	if err != nil {
		return respondWithError(c, err.Error(), fasthttp.StatusInternalServerError)
	}

	return respondWithJSON(c, r)

}

func (s *Server) TransactionsHandler(c *routing.Context) error {
	clientID, err := parseClientID(c.Param("id"))
	if err != nil {
		return respondWithError(c, "Invalid client ID", fasthttp.StatusNotFound)
	}

	input, err := parseAndValidateInput(c.Request.Body())
	if err != nil {
		return respondWithError(c, "Invalid request body", fasthttp.StatusUnprocessableEntity)
	}

	r, err := s.service.HandlePostTransactions(clientID, input)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondWithJSON(c, r)
}

func initializeDatabase(ctx context.Context, dbpool *pgxpool.Pool) *pgxpool.Pool {
	var err error
	for i := 0; i < 5; i++ { // Retry up to 5 times
		dbpool, err = pgxpool.New(ctx, os.Getenv("DB_HOSTNAME"))
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database after retries: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

func respondWithError(c *routing.Context, message string, statusCode int) error {
	c.SetStatusCode(statusCode)
	c.SetContentType("application/json; charset=utf8")
	c.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, message)))
	return nil
}

func respondWithJSON(c *routing.Context, data interface{}) error {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		c.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return nil
	}
	c.SetContentType("application/json; charset=utf8")
	c.SetStatusCode(fasthttp.StatusOK)
	c.Write(jsonResponse)
	return nil
}

func handleServiceError(c *routing.Context, err error) error {
	if response, ok := errorResponseMap[err]; ok {
		return respondWithError(c, response.Message, response.StatusCode)
	}
	return respondWithError(c, "Internal Server Error", fasthttp.StatusInternalServerError)
}
