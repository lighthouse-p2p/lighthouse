package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/lighthouse-p2p/lighthouse/internal/models"
	"github.com/lighthouse-p2p/lighthouse/internal/rtc"
	"github.com/lighthouse-p2p/lighthouse/internal/state"
)

// ProxyHandler creates a proxy over RTC with dynamic connections
type ProxyHandler struct {
	metadata *models.Metadata
	sessions *rtc.Sessions
	st       *state.State
}

// Init initializes the ProxyHandler
func (p *ProxyHandler) Init(metadata *models.Metadata, sessions *rtc.Sessions, st *state.State) {
	p.metadata = metadata
	p.sessions = sessions
	p.st = st
}

// Handler is the fiber "handler"
func (p *ProxyHandler) Handler(ctx *fiber.Ctx) error {
	nickname := ctx.Params("nick")
	if nickname == "" {
		return ctx.SendStatus(400)
	}

	if port, ok := p.sessions.PortMap[nickname]; ok {
		if err := proxy.Do(ctx, fmt.Sprintf("http://localhost:%d", port)); err != nil {
			return err
		}
	} else {
		totalSessions := len(p.sessions.PortMap)
		port := 42001 + totalSessions

		if totalSessions == 9 {
			return ctx.Status(409).SendString("Port overflow")
		}

		newSession := &rtc.Session{}
		newSession.Init(nickname, *p.st, port)

		p.sessions.PortMap[nickname] = port
		p.sessions.RTCSessions[nickname] = newSession

		if err := proxy.Do(ctx, fmt.Sprintf("http://localhost:%d", port)); err != nil {
			return err
		}
	}

	return ctx.SendStatus(404)
}
