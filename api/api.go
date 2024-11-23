package api

import (
	"fmt"

	"github.com/MaxwellGroup/ragexp1/internal/pkg/rag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// APIVersion is the version of the API
	APIVersion = "v1"
)

func NewAPI(e *echo.Echo, rag *rag.Rag) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// put rag pointer in context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("Rag", rag)
			return next(c)
		}
	})

	api := e.Group(fmt.Sprintf("/api/%s", APIVersion))

	api.POST("/upload", UploadHandler)
	api.POST("/ask", AskDocHandler)
	api.GET("/docs", ListAllDocsHandler)
	api.GET("/doc/:id", GetDocHandler)
	api.DELETE("/doc/:id", DeleteDocHandler)
}
