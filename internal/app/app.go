package app

import (
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/theoriginalstove/starshine/internal/fakestrip"
	"github.com/theoriginalstove/starshine/internal/logger"
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
	env     string
	light   lighter
	addr    string
}

type AppOptFunc func(*App) error

func NewApp(opts ...AppOptFunc) (*App, error) {
	a := &App{
		name: "starshine-server",
		addr: "0.0.0.0:7000",
	}
	logger := slog.New(logger.NewHandler(nil))
	a.Logger = logger
	slog.SetDefault(logger)

	a.light = &fakestrip.Fakestrip{}

	for _, fn := range opts {
		if err := fn(a); err != nil {
			return nil, err
		}
	}

	handler := &Handler{
		led: a.light,
	}

	routes := map[string]chi.Router{
		"/": Routes(handler),
	}
	r := NewRouter(routes)

	srv, err := server.NewServer(
		a.name,
		server.WithHandler(r),
		server.WithAddr(a.addr),
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

func WithEnv(e string) AppOptFunc {
	return func(a *App) error {
		a.env = e
		return nil
	}
}

func WithLighter(l lighter) AppOptFunc {
	return func(a *App) error {
		slog.Warn("using lighter", slog.Any("lighter", l))
		a.light = l
		return nil
	}
}

func WithAddr(addr string) AppOptFunc {
	return func(a *App) error {
		a.addr = addr
		return nil
	}
}
