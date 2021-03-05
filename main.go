package main

import (
	"fmt"
	"net/http"
	"vault-seal-watcher/config"
	"vault-seal-watcher/logger"

	"github.com/hashicorp/vault/api"
	"github.com/labstack/echo/v4"
)

// Server context
type Server struct {
	Cfg    *config.Config
	Router *echo.Echo
	Vault  *api.Client
}

// NewServer constructor
func NewServer() *Server {
	cfg := config.Cfg
	vault, err := getVaultClient(cfg.VaultAddr, cfg.VaultTimeout)
	if err != nil {
		logger.Log.Fatal(err)
	}

	return &Server{
		Cfg:    cfg,
		Router: getRouter(),
		Vault:  vault,
	}
}

// RunWebServer API
func (s *Server) RunWebServer() {
	srv := s.Router
	srv.Server.Addr = fmt.Sprintf("%s:%s", s.Cfg.ServerHost, s.Cfg.ServerPort)

	logger.Log.Infof("Starting web server on port %s...", s.Cfg.ServerPort)
	srv.Logger.Fatal(srv.StartServer(srv.Server))
}

func getRouter() *echo.Echo {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, `{"health":"ok"}`)
	})

	return e
}

func main() {
	srv := NewServer()
	go srv.RunWebServer()
	srv.RunVaultWatcher()
}
