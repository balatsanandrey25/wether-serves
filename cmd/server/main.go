package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

const PORT = ":3000"

func main() {
	// Создаем контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инициализируем планировщик
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// останавливаем планировщик при завершении
		if err := s.Shutdown(); err != nil {
			log.Printf("Error shutting down scheduler: %v", err)
		}
	}()

	// Инициализируем cron jobs
	jobs, err := initCron(s)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем планировщик
	s.Start()
	log.Printf("Starting job Id: %v\n", jobs[0].ID())

	// Настраиваем HTTP сервер
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	server := &http.Server{
		Addr:    PORT,
		Handler: r,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Starting server on port: %s", PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ждем сигнал завершения
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Application stopped gracefully")
}

func initCron(scheduler gocron.Scheduler) ([]gocron.Job, error) {
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			1*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("subJob-1_1")
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return []gocron.Job{j}, nil
}
