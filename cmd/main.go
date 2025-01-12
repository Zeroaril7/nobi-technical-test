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

	"github.com/Zeroaril7/nobi-technical-test/config"
	ethereumfetcher "github.com/Zeroaril7/nobi-technical-test/internal/ethereum-fetcher"
	database "github.com/Zeroaril7/nobi-technical-test/pkg/databases"
	redissdk "github.com/Zeroaril7/nobi-technical-test/pkg/redis-sdk"
	"github.com/Zeroaril7/nobi-technical-test/router"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	redissdk.InitRedis()

	database.Connect()

	go scheduler()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	router.SetupRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Config("APP_PORT")),
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %s", config.Config("APP_PORT"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully")
}

func scheduler() {
	c := cron.New()

	_, err := c.AddFunc("@every 1s", func() {
		rate, err := ethereumfetcher.FetchExchangeRate()
		if err != nil {
			log.Println("Error fetching exchange rate:", err)
			return
		}

		err = ethereumfetcher.SaveExchangeRateIntoRedis(rate)
		if err != nil {
			log.Println("Error saving to Redis:", err)
		} else {
			log.Println("Successfully saved exchange rate to Redis:", rate.Text('f', 18))
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	log.Println("Starting scheduler...")
	c.Run()
}
