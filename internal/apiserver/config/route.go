package config

import (
	"github.com/go-fuego/fuego"
	"github.com/woxQAQ/gim/internal/apiserver/controllers"
	"github.com/woxQAQ/gim/internal/apiserver/services"
	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"gorm.io/gorm"
)

func Register(sv *fuego.Server, db *gorm.DB) {
	ustore := stores.NewUserStore(db)
	us := services.NewUserService(ustore)
	uc := controllers.NewUserController(us)
	uc.Route(sv)
}
