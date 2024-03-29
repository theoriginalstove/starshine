package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/theoriginalstove/starshine/internal/app"
	"github.com/theoriginalstove/starshine/internal/lightstrip"
)

var env = os.Getenv("RUN_MODE")

func main() {
	slog.Info("starting up starshine app")
	if env == "" {
		env = "dev"
	}
	opts := []app.AppOptFunc{
		app.WithName(fmt.Sprintf("starshine-%s-apiserver", env)),
	}
	if env == "prod" {
		slog.Info("env is", slog.String("env", env))
		ls := &lightstrip.Lightstrip{}
		err := ls.Init(true)
		if err != nil {
			slog.Error("unable to run Init for lightstrip")
			os.Exit(1)
		}
		opts = append(opts, app.WithLighter(ls))
	}
	app, err := app.NewApp(
		opts...,
	)
	if err != nil {
		slog.Error("unable to create new api app", slog.Any("error", err))
		os.Exit(1)
	}

	app.Server.ListenAndServe()
}
