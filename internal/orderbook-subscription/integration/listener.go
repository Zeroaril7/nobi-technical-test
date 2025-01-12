package integration

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Zeroaril7/nobi-technical-test/config"
	"github.com/Zeroaril7/nobi-technical-test/internal/orderbook-subscription/crypto"
	"github.com/lib/pq"
)

func listenForNewPairs(service crypto.CryptoService, ctx context.Context) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_HOST"), config.Config("DB_PORT"), config.Config("DB_NAME"))

	listener := initPostgresListener(dsn)
	defer listener.Close()

	log.Println("Listening for new pairs via PostgreSQL...")

	for {
		select {
		case notification := <-listener.Notify:
			if notification == nil {
				log.Println("Received nil notification. Skipping...")
				continue
			}

			log.Printf("Notification received: Channel=%s, Extra=%s", notification.Channel, notification.Extra)

			if notification.Extra == "" {
				log.Println("Received notification with empty Extra. Treating as DELETE operation.")
				handleDeleteNotification(service, ctx)
			} else {
				handleNotification(service, ctx, notification.Extra)
			}

		case <-time.After(30 * time.Second):
			go checkPostgresListenerConnection(&listener, dsn)
		}
	}
}

func handleDeleteNotification(service crypto.CryptoService, ctx context.Context) {
	log.Println("Handling delete notification: Fetching updated pairs...")

	newPairs, err := fetchInitialPairs(service, ctx)
	if err != nil {
		log.Printf("Failed to fetch updated pairs: %v", err)
		return
	}

	if hasPairChanges(newPairs) {
		mu.Lock()
		activePairs = newPairs
		mu.Unlock()

		log.Println("Pairs updated after delete operation. Triggering reconnect...")
		reconnectChan <- true
	}
}

func initPostgresListener(dsn string) *pq.Listener {
	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, func(event pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("PostgreSQL listener error: %v", err)
		}
	})

	listener.Notify = make(chan *pq.Notification, 10)

	if err := listener.Listen("new_crypto"); err != nil {
		log.Fatalf("Failed to start PostgreSQL listener: %v", err)
	}

	return listener
}

func handleNotification(service crypto.CryptoService, ctx context.Context, extra string) {
	log.Printf("New crypto pair notification received: %s", extra)

	newPairs, err := fetchInitialPairs(service, ctx)
	if err != nil {
		log.Printf("Failed to fetch updated pairs. Retrying... Error: %v", err)
		time.Sleep(5 * time.Second)
		return
	}

	if hasPairChanges(newPairs) {
		mu.Lock()
		activePairs = newPairs
		mu.Unlock()

		log.Println("Pairs updated. Triggering reconnect...")
		reconnectChan <- true
	}
}

func hasPairChanges(newPairs []string) bool {
	mu.Lock()
	defer mu.Unlock()

	if len(activePairs) != len(newPairs) {
		return true
	}

	for i, pair := range activePairs {
		if pair != newPairs[i] {
			return true
		}
	}

	return false
}

func checkPostgresListenerConnection(listener **pq.Listener, dsn string) {
	log.Println("Checking PostgreSQL listener connection...")
	if err := (*listener).Ping(); err != nil {
		log.Printf("PostgreSQL listener ping failed: %v. Reconnecting...", err)

		(*listener).Close()

		newListener := initPostgresListener(dsn)
		*listener = newListener

		log.Println("PostgreSQL listener reconnected successfully.")
	}
}
