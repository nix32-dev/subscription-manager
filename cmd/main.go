package main

import (
	"context"
	"net/http"
	"os/signal"
	_ "subscription/docs"
	"subscription/internal/handler"
	"subscription/internal/repository"
	"subscription/internal/service"
	"syscall"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Subscription Manager
// @version 1.0
// @description API Server for Subscription Manager Application

// @contact.name   Kirill | nix32
// @contact.email  nix32.dev@gmail.com

// @host localhost:8080
// @BasePath /

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	service.CTX = &ctx

	logFile, logger := service.InitLogger()
	defer logFile.Close()
	defer logger.Sync()

	if err := repository.CreateDBConnection(*service.CTX); err != nil {
		service.Logger.Fatal(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/create", handler.CreateSubscriptionH)
	mux.HandleFunc("/change", handler.ChangeSubscriptionH)
	mux.HandleFunc("/delete", handler.DeleteSubscriptionH)
	mux.HandleFunc("/list", handler.GetSubscriptionH)
	mux.Handle("/swagger/", httpSwagger.Handler())

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	stop := *service.CTX
	go func() {
		if err := server.ListenAndServe(); err != nil {
			service.Logger.Fatal(err.Error())
		}
		<-stop.Done()
		if err := server.Shutdown(*service.CTX); err != nil {
		}
	}()
	<-stop.Done()
}
