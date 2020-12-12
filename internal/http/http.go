package http

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
	"github.com/lighthouse-p2p/lighthouse/internal/api"
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
	time.Sleep(500 * time.Millisecond)
	engine := html.New("./homepage", ".html")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
		DisableKeepalive:      true,
	})

	app.Get("/coins", func(c *fiber.Ctx) error {
		coins, err := api.Coins(fmt.Sprintf("http://%s/v1/coins", metadata.Host), metadata.PubKey)
		if err != nil {
			return err
		}

		return c.SendString(coins)
	})

	app.Get("/stats", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("Number of Goroutines: %d", runtime.NumGoroutine()))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"PubKey":   metadata.PubKey,
			"NickName": metadata.NickName,
		})
	})
	app.Static("/", "./homepage")

	proxyHandler := handlers.ProxyHandler{}
	proxyHandler.Init(metadata, sessions, st)

	app.Use("/proxy/:nick/*", proxyHandler.Handler)

	fmt.Printf("  %s\n", aurora.Bold(aurora.Green("Proxy server is up ✓")))

	go log.Fatal(app.Listen("127.0.0.1:42000"))
}
