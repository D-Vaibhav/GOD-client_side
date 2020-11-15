package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/vaibhav/GOD/client_side/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// specifying what tool to be used to log
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)

	// creating product handler
	// this is the advantage of abstracting we have to create product handler and then only we have it's methods
	productHandler := handlers.NewProducts(l)

	// creating root router
	rootRouter := mux.NewRouter() // this is same as we done with our product handler

	// SUBROUTERS out of the rootRouter
	getRouter := rootRouter.Methods(http.MethodGet).Subrouter()
	postRouter := rootRouter.Methods(http.MethodPost).Subrouter()
	putRouter := rootRouter.Methods(http.MethodPut).Subrouter()

	// MIDDLEWARES
	postRouter.Use(productHandler.MiddlewareValidateProduct)
	putRouter.Use(productHandler.MiddlewareValidateProduct) // lacks complete logic

	// HANDLEFUCN
	// router specific requests handling
	getRouter.HandleFunc("/", productHandler.GetProducts)
	postRouter.HandleFunc("/", productHandler.AddProduct)
	putRouter.HandleFunc("/{id:[0-9]+}", productHandler.UpdateProduct)

	// SERVER
	clientServer := http.Server{
		Addr:         ":8000",           // configure this bindAddress
		Handler:      rootRouter,        // set the default handler
		ErrorLog:     l,                 // set the Logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from client
		WriteTimeout: 10 * time.Second,  // max time to write request to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP keep-alive
	}

	// starting server
	go func() {
		l.Println("Starting server...")

		// gorilla provide ListenAndServe() which will do both listening then serve
		err := clientServer.ListenAndServe()
		if err != nil {
			l.Fatalf("Failed to start server, ERROR:", err)
		}
		l.Println("Server started on port:8000")
	}()

	// trapping signals like Interrupt and kill, before shut-down
	signalNotificationChannel := make(chan os.Signal)
	signal.Notify(signalNotificationChannel, os.Interrupt) // if interrupted
	signal.Notify(signalNotificationChannel, os.Kill)      // if killed

	sig := <-signalNotificationChannel
	l.Println("signal is notified by", sig)

	timeout := 30 * time.Second

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err := clientServer.Shutdown(ctx)
	if err != nil {
		l.Fatalf("Failed to shutdown server, ERROR: %s", err)
	}
}
