package database

import (
	"log"

	"gorm.io/gorm"
)

func createTrigger(db *gorm.DB) error {

	createFunctionSQL := `
	CREATE OR REPLACE FUNCTION notify_new_crypto()
	RETURNS TRIGGER AS $$
	BEGIN
	    PERFORM pg_notify('new_crypto', NEW.pair);
	    RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`

	checkInsertTriggerSQL := `
	SELECT COUNT(1)
	FROM pg_trigger
	WHERE tgname = 'trigger_new_crypto_insert';
	`

	checkDeleteTriggerSQL := `
	SELECT COUNT(1)
	FROM pg_trigger
	WHERE tgname = 'trigger_new_crypto_delete';
	`

	createInsertTriggerSQL := `
	CREATE TRIGGER trigger_new_crypto_insert
	AFTER INSERT ON crypto
	FOR EACH ROW
	EXECUTE FUNCTION notify_new_crypto();
	`

	createDeleteTriggerSQL := `
	CREATE TRIGGER trigger_new_crypto_delete
	AFTER DELETE ON crypto
	FOR EACH ROW
	EXECUTE FUNCTION notify_new_crypto();
	`

	log.Println("Creating function notify_new_crypto if not exists...")
	if err := db.Exec(createFunctionSQL).Error; err != nil {
		log.Printf("Failed to create function notify_new_crypto: %v", err)
		return err
	}
	log.Println("Function notify_new_crypto created successfully.")

	var insertTriggerCount int
	log.Println("Checking existence of INSERT trigger...")
	if err := db.Raw(checkInsertTriggerSQL).Scan(&insertTriggerCount).Error; err != nil {
		log.Printf("Failed to check INSERT trigger existence: %v", err)
		return err
	}
	log.Printf("INSERT trigger count: %d", insertTriggerCount)

	if insertTriggerCount == 0 {
		log.Println("Creating INSERT trigger trigger_new_crypto_insert...")
		if err := db.Exec(createInsertTriggerSQL).Error; err != nil {
			log.Printf("Failed to create INSERT trigger trigger_new_crypto_insert: %v", err)
			return err
		}
		log.Println("INSERT trigger trigger_new_crypto_insert created successfully.")
	} else {
		log.Println("INSERT trigger trigger_new_crypto_insert already exists. Skipping creation.")
	}

	var deleteTriggerCount int
	log.Println("Checking existence of DELETE trigger...")
	if err := db.Raw(checkDeleteTriggerSQL).Scan(&deleteTriggerCount).Error; err != nil {
		log.Printf("Failed to check DELETE trigger existence: %v", err)
		return err
	}
	log.Printf("DELETE trigger count: %d", deleteTriggerCount)

	if deleteTriggerCount == 0 {
		log.Println("Creating DELETE trigger trigger_new_crypto_delete...")
		if err := db.Exec(createDeleteTriggerSQL).Error; err != nil {
			log.Printf("Failed to create DELETE trigger trigger_new_crypto_delete: %v", err)
			return err
		}
		log.Println("DELETE trigger trigger_new_crypto_delete created successfully.")
	} else {
		log.Println("DELETE trigger trigger_new_crypto_delete already exists. Skipping creation.")
	}

	return nil
}
