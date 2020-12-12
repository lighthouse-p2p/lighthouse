package http

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/lighthouse-p2p/lighthouse/internal/http/handlers"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/lighthouse-p2p/lighthouse/internal/rtc"
	"github.com/lighthouse-p2p/lighthouse/internal/state"
	"github.com/logrusorgru/aurora"
)

// InitFileServer initializes the data server
func InitFileServer(metadata models.Metadata) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		DisableKeepalive:      true,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Request().Header.Add("Cache-Control", "no-cache")
		return c.Next()
	})

	app.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("./data"),
		Browse: true,
	}))

	fmt.Printf("  %s\n", aurora.Bold(aurora.Green("Data server is up ✓")))

	go log.Fatal(app.Listen("127.0.0.1:42011"))
}

// InitProxyServer initializes the proxy server
func InitProxyServer(metadata *models.Metadata, sessions *rtc.Sessions, st *state.State) {
	time.Sleep(time.Second)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/stats", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("Number of Goroutines: %d", runtime.NumGoroutine()))
	})

	proxyHandler := handlers.ProxyHandler{}
	proxyHandler.Init(metadata, sessions, st)

	app.Use("/proxy/:nick/*", proxyHandler.Handler)

	fmt.Printf("  %s\n", aurora.Bold(aurora.Green("Proxy server is up ✓")))

	go log.Fatal(app.Listen("127.0.0.1:42000"))
}
