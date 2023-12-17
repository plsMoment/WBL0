package main

import (
	"WBL0/config"
	"WBL0/internal/db"
	"WBL0/internal/model"
	streamnats "WBL0/internal/nats-streaming"
	"WBL0/migrations"
	orderservice "WBL0/service"
	"WBL0/transport"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	clusterID = "test-cluster"
	clientID  = "consumer"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(fmt.Errorf("cannot parse config: %w", err))
	}

	// Init connections pool
	store, err := db.New(context.Background(), cfg)
	if err != nil {
		os.Exit(1)
	}
	defer store.Close()

	if err = migrations.MigrateUp(cfg); err != nil {
		log.Printf("migrate up failed: %v", err)
		return
	}

	// connect to STAN cluster
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Printf("creating connection to cluster failed: %v", err)
	}
	defer func() {
		if err := sc.Close(); err != nil {
			log.Printf("closing connection to cluster failed: %v", err)
		}
	}()

	orderCh := make(chan model.Order, 20)
	defer close(orderCh)

	// Create new subscriber
	sub, err := streamnats.NewSubscriber(sc, orderCh)
	if err != nil {
		log.Printf("creating new subscriber failed: %v", err)
	}
	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("unsubscribe failed: %v", err)
		}
	}()

	service := orderservice.New(store, orderCh)

	// Init router
	router := chi.NewRouter()
	router.Get("/", transport.HomePage())
	router.Get("/order/{order_id}", transport.GetOrder(service))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.ServerHostAddress,
		Handler: router,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("Server stopped")
			} else {
				log.Printf("error during server shutdown: %v", err)
			}
		}
	}()

	log.Println("Server started")

	<-done
	log.Println("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("failed to stop server: %v", err)
		return
	}
}
