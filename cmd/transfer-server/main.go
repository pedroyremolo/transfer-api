package main

import (
	"context"
	"fmt"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var adder adding.Service
	var port int

	storage, err := mongodb.NewStorageFromEnv()
	if err != nil {
		panic(err.Error())
	}

	dbCtx := context.Background()
	storage.Connect(dbCtx)
	storage.CreateIndexes(dbCtx)
	defer storage.Disconnect(dbCtx)

	adder = adding.NewService(storage)

	handler := rest.Handler(adder)
	port, err = strconv.Atoi(os.Getenv("APP_PORT"))
	portStr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server at port %s", portStr)
	log.Fatal(http.ListenAndServe(portStr, handler))
}
