package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "effective-mobile-task/docs"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"effective-mobile-task/internal/client"
	"effective-mobile-task/internal/config"
	"effective-mobile-task/internal/handler"
	"effective-mobile-task/internal/repository"
	"effective-mobile-task/internal/service"
)

// @title People Enrichment API
// @version 1.0
// @description Сервис для обогащения ФИО возрастом, полом и национальностью
// @host localhost:8080
// @BasePath /

// @contact.name API Support
// @contact.email support@example.com
func main() {
	cfg := config.LoadConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting server", zap.String("host", cfg.ServerHost), zap.String("port", cfg.ServerPort))

	dsn := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser +
		" password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " sslmode=" + cfg.DBSSLMode

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("Cannot connect to DB", zap.Error(err))
	}
	defer db.Close()

	enricher := client.NewEnrichmentClient(cfg.GenderizeAPIURL, cfg.AgifyAPIURL, cfg.NationalizeAPIURL)
	repo := repository.NewPersonRepository(db)
	service := service.NewPersonService(repo, enricher)
	handler := handler.NewPersonHandler(service, logger)

	r := mux.NewRouter()
	r.HandleFunc("/persons", handler.Create).Methods("POST")
	r.HandleFunc("/persons", handler.List).Methods("GET")
	r.HandleFunc("/persons/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/persons/{id}", handler.Update).Methods("PUT")
	r.HandleFunc("/persons/{id}", handler.Delete).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:         cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server stopped")
}
