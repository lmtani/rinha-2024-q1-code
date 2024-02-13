package main

import (
	"encoding/json"
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

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
