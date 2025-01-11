package main

import (
	"fmt"
	"log"

	"github.com/Zeroaril7/nobi-technical-test/config"
	ethereumfetcher "github.com/Zeroaril7/nobi-technical-test/internal/ethereum-fetcher"
	redissdk "github.com/Zeroaril7/nobi-technical-test/pkg/redis-sdk"
	"github.com/Zeroaril7/nobi-technical-test/router"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	redissdk.InitRedis()

	go scheduler()

	r := gin.Default()
	router.SetupRoutes(r)

	log.Printf("Starting server on port %s", config.Config("APP_PORT"))
	err := r.Run(fmt.Sprintf(":%s", config.Config("APP_PORT")))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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
