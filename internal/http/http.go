package http

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/logrusorgru/aurora"
)

// InitFileServer initializes the data server
func InitFileServer(metadata models.Metadata) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/stats", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("Number of Goroutines: %d", runtime.NumGoroutine()))
	})

	app.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("./data"),
		Browse: true,
	}))

	fmt.Printf("  %s\n", aurora.Bold(aurora.Green("Data server is up ✓")))

	go log.Fatal(app.Listen("127.0.0.1:42011"))
}
