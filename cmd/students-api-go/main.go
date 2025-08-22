package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saurav-lal-karn/students-api-go/internal/config"
	"github.com/saurav-lal-karn/students-api-go/internal/http/handlers/student"
	"github.com/saurav-lal-karn/students-api-go/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	// Setup loggers
	// Setup database
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	// Setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started:", slog.String("address", cfg.Addr))
	// fmt.Printf("Server started: %s", cfg.HTTPServer.Addr)

	// Run in goroutine for the graceful shutdown
	// This is must in production
	// otherwise ongoing request will be stopped immediately
	// This is not preferred in production
	done := make(chan os.Signal, 1)

	// Get the interrupt signal from os as well as others as well
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	// Wait for signal to arrive in done channel
	// Until then, it will not move forward
	<-done

	// Use structured log over traditional log
	slog.Info("Shutting Down the Server")

	// We can shutdown immediately
	// But sometimes the system gets in infinite loop
	// This will cause the system to never shut down
	// So use timer to know if the server has stopped over certain time or not
	// We need to do graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Above context waits for 5 seconds as mentioned in timeout
	// If not, we will shutdown the server
	// err := server.Shutdown(ctx)
	// if err != nil {
	// 	slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	// }

	// Shorthand way to write above error condition
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
