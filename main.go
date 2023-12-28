package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/dotenv-org/godotenvvault"

	"github.com/Kbgjtn/notethingness-api.git/api"
)

func initServer(ctx context.Context) error {
	if err := godotenvvault.Load("./.env.local"); err != nil {
		return fmt.Errorf(" [ 💢Cannot load .env.local! ] " + "\nError: " + err.Error())
	}
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
