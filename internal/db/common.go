package db

import (
	"fmt"
	"log"

	vad "github.com/asaskevich/govalidator"
	"github.com/woxQAQ/gim/internal/models"
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

func userStructValid(user models.UserBasic) (err error) {
	if ok, err := vad.ValidateStruct(user); !ok {
		return fmt.Errorf("用户结构不完整")
	} else if err != nil {
		return err
	}
	return
}
