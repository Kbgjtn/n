package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/Kbgjtn/notethingness-api.git/api"
)

func initServer(ctx context.Context) error {
	server := api.NewServer()
	return server.Start(ctx)
}

func main() {
	serverCtx, stopCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopCtx()

	if err := initServer(serverCtx); err != nil {
		log.Fatal("can't run the server: ", err)
		slog.Error(" [ 💢Cannot run the server! ] " + "\nError: " + err.Error())
		os.Exit(1)
	}
}
