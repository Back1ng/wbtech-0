package main

import (
	"context"
	"errors"
	"github.com/Back1ng/wbtech-0/internal/cache"
	natsclient "github.com/Back1ng/wbtech-0/internal/nats"
	"github.com/Back1ng/wbtech-0/internal/repository"
	"github.com/Back1ng/wbtech-0/internal/rest"
	"github.com/Back1ng/wbtech-0/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Back1ng/wbtech-0/internal/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	connOpts := postgres.SetupOptions{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	// Execute migrations
	{
		m, err := migrate.New(
			"file://migrations",
			connOpts.String(),
		)
		if err != nil {
			log.Fatal("Error when init migrations:", err)
		}

		if err := m.Up(); err != nil {
			if !errors.Is(migrate.ErrNoChange, err) {
				log.Fatal("Error when apply migrations:", err)
			}
		}
	}

	// get connection string
	pool, err := postgres.New(ctx, connOpts)
	if err != nil {
		log.Fatal("Error when postgres.New:", err)
	}

	repo := repository.NewOrdersRepo(pool)
	memoryCache := cache.New()
	uc := usecase.NewOrderUsecase(pool, repo, memoryCache)

	// Move orders to cache
	{
		orders, err := uc.GetAllOrders(ctx)
		if err != nil {
			log.Fatal("Error when get all orders:", err)
		}

		memoryCache.StoreAll(orders)
	}

	go func() {
		mux := rest.NewHandler(uc)
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Println("http.ListenAndServe", err)
			return
		}
	}()

	nc, err := natsclient.NewClient()
	if err != nil {
		log.Fatal("natsclient.NewClient():", err)
	}
	defer nc.Close()

	log.Println(memoryCache.GetAll())

	orderCh := nc.ListenOrdersSubject("ordersSubject")
	go func() {
		for order := range orderCh {
			if err := uc.StoreOrder(ctx, order); err != nil {
				log.Println("uc.StoreOrder", err)
			}
		}
	}()

	<-done

	log.Println("Shutting down...")
	close(orderCh)
}
