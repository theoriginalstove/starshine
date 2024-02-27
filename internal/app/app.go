package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/theoriginalstove/starshine/internal/server"
)

var publicRoutes map[string]chi.Router

func init() {
}

type App struct {
	name    string
	Server  *server.Server
	Routers []chi.Router
	Logger  *slog.Logger
}

type AppOptFunc func(*App) error

func NewApp(opts ...AppOptFunc) (*App, error) {
	a := &App{
		name: "sharedo-api-service",
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	a.Logger = logger
	slog.SetDefault(logger)

	for _, fn := range opts {
		if err := fn(a); err != nil {
			return nil, err
		}
	}

	routes := map[string]chi.Router{
		"/": Routes(),
	}
	r := NewRouter(routes)

	srv, err := server.NewServer(
		a.name,
		server.WithHandler(r),
		server.WithAddr("0.0.0.0:8080"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create new server for app: %w", err)
	}
	slog.Info("created new server with address", slog.Any("server-addr", srv.Addr()))

	a.Server = srv

	return a, nil
}

func WithName(name string) AppOptFunc {
	return func(a *App) error {
		a.name = name
		return nil
	}
}

func WithLogger(l *slog.Logger) AppOptFunc {
	return func(a *App) error {
		a.Logger = l
		slog.SetDefault(l)
		return nil
	}
}
