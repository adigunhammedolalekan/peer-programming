package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	userStore := NewUserStore()
	transactionStore := NewTransactionStore()

	processor := NewProcessor(userStore, transactionStore, workerCountFromConfig())
	processor.Start()

	apiService := NewApiService(userStore, transactionStore, processor)

	apiHandler := NewHandler(apiService)

	router := chi.NewRouter()

	router.Post("/users", apiHandler.CreateUser)
	router.Post("/transactions", apiHandler.CreateTransaction)
	router.Get("/users", apiHandler.GetUsers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}
	port = fmt.Sprintf(":%s", port)

	go func() {
		log.Printf("Server started %s", port)
		if err := http.ListenAndServe(port, router); err != nil {
			log.Fatal("failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	processor.Stop()
}

const DefaultWorkerCount = 5

func workerCountFromConfig() int {
	count, err := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if err != nil || count <= 0 {
		return DefaultWorkerCount
	}
	return count
}
