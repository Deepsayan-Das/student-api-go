package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Deepsayan-Das/student-api-go/internal/config"
	"github.com/Deepsayan-Das/student-api-go/internal/http/handlers/student"
	"github.com/Deepsayan-Das/student-api-go/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()
	//use logger
	//setup db
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage %s", err.Error())
	}
	slog.Info("storage initialized successfully", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	//setup server

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}
	fmt.Printf("\n server Started %s", cfg.HTTPServer.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatalf("Failed to sstart server %s", err.Error())
		}
	}()

	//Graceful shutdown
	<-done
	slog.Info("shutting the server down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown sucessfull")
}
