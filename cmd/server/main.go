package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/theoriginalstove/starshine/internal/app"
)

var env = os.Getenv("RUN_MODE")

func main() {
	if env == "" {
		env = "dev"
	}
	app, err := app.NewApp(
		app.WithName(fmt.Sprintf("starshine-%s-apiserver", env)),
		// app.WithConnection(context.Background(), config.Secrets()["ROACH_CONN"]),
	)
	if err != nil {
		slog.Error("unable to create new api app", slog.Any("error", err))
		os.Exit(1)
	}

	app.Server.ListenAndServe()
}
