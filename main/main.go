package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	validate "github.com/MohamedNazir/Validate"
	//"validate"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	// our validate service
	srv := validate.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := validate.Endpoints{
		GetEndpoint:      validate.MakeGetEndpoint(srv),
		StatusEndpoint:   validate.MakeStatusEndpoint(srv),
		ValidateEndpoint: validate.MakeValidateEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("validate is listening on port:", *httpAddr)
		handler := validate.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
