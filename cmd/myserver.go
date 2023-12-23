package cmd

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MyServer struct{}

func (s *MyServer) PostArticlesStart(ctx echo.Context) error {
	// Implement your logic here
	return ctx.String(http.StatusOK, "Success")
}

func (s *MyServer) GetPing(ctx echo.Context) error {
	// Implement your logic here
	return ctx.String(http.StatusOK, "Pong")
}
