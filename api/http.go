package api

import (
	"fmt"

	"github.com/elchemista/easy_rag/internal/pkg/rag"
	"github.com/labstack/echo/v4"
)

const (
	// APIVersion is the version of the API
	APIVersion = "v1"
)

func StartServer(e *echo.Echo, rag *rag.Rag) {
	e.POST(fmt.Sprintf("/api/%s/upload", APIVersion), func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.POST(fmt.Sprintf("/api/%s/ask", APIVersion), func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	// result, err := rag.Database.ListDocuments()

	// log.Printf("result: %v", result)
	// log.Println(err)

	e.Logger.Fatal(e.Start(":4002"))
}
