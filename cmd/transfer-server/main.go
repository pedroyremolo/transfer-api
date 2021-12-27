package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/gatekeeper/jwt"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
	ah "github.com/pedroyremolo/transfer-api/pkg/http/rest/adding"
	auh "github.com/pedroyremolo/transfer-api/pkg/http/rest/authenticating"
	lh "github.com/pedroyremolo/transfer-api/pkg/http/rest/listing"
	th "github.com/pedroyremolo/transfer-api/pkg/http/rest/transferring"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"github.com/pedroyremolo/transfer-api/pkg/transferring"
	"github.com/pedroyremolo/transfer-api/pkg/updating"
	"github.com/sirupsen/logrus"
)

func main() {
	log := lgr.NewDefaultLogger()
	logger := logrus.NewEntry(log)

	var port int

	storage, err := mongodb.NewStorageFromEnv()
	if err != nil {
		logger.Fatalf("failed to get storage: %s", err)
	}

	gatekeeper := jwt.NewGatekeeperFromEnv()

	dbCtx := context.Background()
	storage.Connect(dbCtx)
	storage.CreateIndexes(dbCtx)
	defer storage.Disconnect(dbCtx)

	adder := adding.NewService(storage)
	lister := listing.NewService(storage)
	authenticator := authenticating.NewService(storage, gatekeeper)
	transferor := transferring.NewService()
	updater := updating.NewService(storage)

	addingHandler := ah.NewHandler(logger, adder)
	transferringHandler := th.NewHandler(logger, transferor, lister, adder, updater)
	listingHandler := lh.NewHandler(logger, lister)
	authenticatingHandler := auh.NewHandler(logger, authenticator, lister)

	handler := rest.Handler(logger, addingHandler, transferringHandler, authenticatingHandler, listingHandler)
	port, err = strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8080
	}
	portStr := fmt.Sprintf(":%d", port)
	log.Infof("Starting server at port %s", portStr)
	log.Fatal(http.ListenAndServe(portStr, handler))
}
