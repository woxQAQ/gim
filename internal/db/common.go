package db

import (
	"log"

	"gorm.io/gorm"
)

func closeTransactions(tx *gorm.DB, err error) {
	if err != nil {
		log.Println("Error occurred:", err)
		tx.Rollback()
		return
	} else {
		err = tx.Commit().Error
		if err != nil {
			log.Println("Failed to commit transaction:", err)
			tx.Rollback()
		}
	}
}
