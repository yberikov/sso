package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sso/internal/config"
	"sso/internal/http-server/handlers/editUser"
	"sso/internal/http-server/handlers/login"
	"sso/internal/http-server/handlers/register"
	middlewareCustom "sso/internal/http-server/middleware"
	"sso/internal/storage/sqlite"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Fatalln(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/register", register.New(storage))
	r.Post("/login", login.New(storage))

	r.Group(func(r chi.Router) {
		r.Use(middlewareCustom.AuthMiddleware)
		r.Post("/edit", editUser.New(storage))
	})

	fmt.Printf("Server started at %v", cfg.Address)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln("failed to start server")
		}
	}()

	log.Println("server started")

	<-done
	log.Println("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("failed to stop server", err)
		return
	}

	if err := storage.Stop(); err != nil {
		log.Fatalln("failed to stop storage", err)
		return
	}

	log.Println("server stopped")
}
