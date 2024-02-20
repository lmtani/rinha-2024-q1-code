package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lmtani/rinha-2024-q1-code/internal/repositories"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmtani/rinha-2024-q1-code/internal/services"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type Server struct {
	dbpool  *pgxpool.Pool
	service *services.Service
}

func NewServer(dbpool *pgxpool.Pool) *Server {
	repository := repositories.NewPostgresRepository(dbpool)
	return &Server{dbpool: dbpool, service: services.NewService(repository)}

}

func main() {
	ctx := context.Background()
	dbpool := initializeDatabase(ctx)
	defer dbpool.Close()
	// GET PORT from env var
	port := os.Getenv("PORT")

	server := NewServer(dbpool)

	fmt.Println(fmt.Sprintf("Server running on port %s", port))

	router := routing.New()
	router.Get("/clientes/<id>/extrato", server.StatementHandler)
	router.Post("/clientes/<id>/transacoes", server.TransactionsHandler)

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), router.HandleRequest))
}

func (s *Server) StatementHandler(c *routing.Context) error {
	clientID, err := parseClientID(c.Param("id"))
	if err != nil {
		return respondWithError(c, "Invalid client ID", fasthttp.StatusNotFound)
	}

	r, err := s.service.HandleGetStatement(clientID)
	if err != nil {
		return handleServiceError(c, err)
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

	r, err := s.service.HandlePostTransactions(clientID, &input)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondWithJSON(c, r)
}

func initializeDatabase(ctx context.Context) *pgxpool.Pool {
	var (
		dbpool *pgxpool.Pool
		err    error
	)
	for i := 0; i < 5; i++ { // Retry up to 5 times
		dbpool, err = pgxpool.New(ctx, os.Getenv("DB_HOSTNAME"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			continue
		}
		err = dbpool.Ping(ctx)
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		time.Sleep(2 * time.Second)
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
