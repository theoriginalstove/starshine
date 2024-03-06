package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/theoriginalstove/starshine/internal/app"
)

var env = os.Getenv("RUN_MODE")

func main() {
	slog.Info("starting up starshine app")
	if env == "" {
		env = "dev"
	}
	opts := []app.AppOptFunc{
		app.WithName(fmt.Sprintf("starshine-%s-apiserver", env)),
		app.WithAddr(":9020"),
	}
	app, err := app.NewApp(
		opts...,
	)
	if err != nil {
		slog.Error("unable to create new api app", slog.Any("error", err))
		os.Exit(1)
	}
	err = app.Server.ListenAndServe()
	slog.Error("error encounted running server", slog.Any("err", err))
}
