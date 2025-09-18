package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fiber-rest-api/internal/db"
	"fiber-rest-api/internal/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// initialize sqlite database (file: data.db)
	if err := db.Init("data.db"); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	app := fiber.New()
	router.SetupRoutes(app)

	// start server in background
	srvErr := make(chan error, 1)
	go func() {
		log.Println("starting server on :3000")
		if err := app.Listen(":3000"); err != nil {
			srvErr <- err
		}
	}()

	// handle shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErr:
		log.Fatalf("server error: %v", err)
	case sig := <-stop:
		log.Printf("received signal %v, shutting down...", sig)

		// allow up to 5s for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			if err := app.Shutdown(); err != nil {
				log.Printf("error during shutdown: %v", err)
			}
			close(done)
		}()

		select {
		case <-done:
			log.Println("server stopped gracefully")
		case <-ctx.Done():
			log.Println("graceful shutdown timed out")
		}
	}
}
