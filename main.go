package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/dotenv-org/godotenvvault"

	"github.com/Kbgjtn/notethingness-api.git/api"
)

func main() {
	if err := godotenvvault.Load(".env.local"); err != nil {
		slog.Error(" [ 💢Cannot load .env file ] ")
		os.Exit(1)
	}

	serverCtx, stopCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopCtx()

	server := api.NewServer()

	if err := server.Start(serverCtx); err != nil {
		log.Fatal("can't run the server: ", err)
		slog.Error(" [ 💢Cannot run the server! ] " + "\nError: " + err.Error())
		os.Exit(1)
	}
}
