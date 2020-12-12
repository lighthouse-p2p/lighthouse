package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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

	validationRegex := regexp.MustCompile("^[a-z]+$")

	if !validationRegex.Match([]byte(nickname)) {
		return errors.New("The nickname is invalid")
	}

	if nickname == p.metadata.NickName {
		return ctx.Status(401).SendString("Sorry can't connect proxy to myself")
	}

	subPath := strings.Split(ctx.Path(), fmt.Sprintf("/%s", nickname))[1]

	if port, ok := p.sessions.PortMap[nickname]; ok {
		return ctx.Redirect(fmt.Sprintf("http://localhost:%d%s", port, subPath))
	}

	totalSessions := len(p.sessions.PortMap)
	port := 42001 + totalSessions

	if totalSessions == 9 {
		return ctx.Status(409).SendString("Port overflow")
	}

	newSession := &rtc.Session{}
	err := newSession.Init(nickname, *p.st, port)

	if err != nil {
		return err
	}

	p.sessions.PortMap[nickname] = port
	p.sessions.RTCSessions[nickname] = newSession

	time.Sleep(1500 * time.Millisecond)

	return ctx.Redirect(fmt.Sprintf("http://localhost:%d%s", port, subPath))
}
