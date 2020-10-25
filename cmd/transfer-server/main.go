package main

import (
	"context"
	"fmt"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/gatekeeper/jwt"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"github.com/pedroyremolo/transfer-api/pkg/transferring"
	"github.com/pedroyremolo/transfer-api/pkg/updating"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var adder adding.Service
	var lister listing.Service
	var port int

	storage, err := mongodb.NewStorageFromEnv()
	if err != nil {
		panic(err.Error())
	}

	gatekeeper := jwt.NewGatekeeperFromEnv()

	dbCtx := context.Background()
	storage.Connect(dbCtx)
	storage.CreateIndexes(dbCtx)
	defer storage.Disconnect(dbCtx)

	adder = adding.NewService(storage)
	lister = listing.NewService(storage)
	authenticator := authenticating.NewService(storage, gatekeeper)
	transferor := transferring.NewService()
	updater := updating.NewService(storage)

	handler := rest.Handler(adder, lister, authenticator, transferor, updater)
	port, err = strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8080
	}
	portStr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server at port %s", portStr)
	log.Fatal(http.ListenAndServe(portStr, handler))
}
