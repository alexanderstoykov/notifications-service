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

	"github.com/alexanderstoykov/notifications-service/config"
	http_server "github.com/alexanderstoykov/notifications-service/internal/api/http"
	"github.com/alexanderstoykov/notifications-service/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	conn, err := storage.NewConnection(cfg.Database)
	if err != nil {
		return err
	}

	_, cancel := context.WithCancel(context.Background())

	notificationRepo := storage.NewNotificationRepository(conn)
	handler := http_server.NewHandler(notificationRepo)

	router := gin.Default()
	http_server.RegisterRoutes(router, handler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errs := make(chan error, 1)
	go func() {
		errs <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		cancel()
	}()

	select {
	case err := <-errs:
		return fmt.Errorf("server error: %w", err)
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout))
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			//nolint: errcheck
			srv.Close()

			return fmt.Errorf("could not stop server gracefully: %v", err)
		}
	}

	return nil
}

//go func() {
//	log.Printf("[Server] listen on %s", server.Addr)
//
//	serverErrors <- server.ListenAndServe()
//}()
//
//shutdown := make(chan os.Signal, 1)
//signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
//
//select {
//case err := <-serverErrors:
//return fmt.Errorf("server error: %w", err)
//case sig := <-shutdown:
//log.Println("shutdown started", sig)
//defer log.Println("shutdown completed", sig)
//
//ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
//defer cancel()
//
//if err := server.Shutdown(ctx); err != nil {
////nolint: errcheck
//server.Close()
//
//return fmt.Errorf("could not stop server gracefully: %v", err)
//}
//}
