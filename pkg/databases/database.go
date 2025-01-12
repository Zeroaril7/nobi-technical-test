package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Zeroaril7/nobi-technical-test/config"
	"github.com/Zeroaril7/nobi-technical-test/internal/orderbook-subscription/crypto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	Instance *gorm.DB
}

var DB DBInstance

func Connect() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_HOST"), config.Config("DB_PORT"), config.Config("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("Database connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	DB = DBInstance{
		Instance: db,
	}

	err = db.AutoMigrate(&crypto.CryptoEntity{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = createTrigger(db)
	if err != nil {
		log.Fatalf("Failed to create trigger: %v", err)
	}
}
