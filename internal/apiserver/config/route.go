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
	mstore := stores.NewMessageStore(db)
	us := services.NewUserService(ustore)
	ms := services.NewMessageService(mstore)
	uc := controllers.NewUserController(us)
	mc := controllers.NewMessageController(ms)
	apiv1 := fuego.Group(sv, "/api/v1",
		fuego.OptionDescription("API v1"),
		fuego.OptionHeader("Authentication", "Bearer Token", fuego.ParamRequired()),
	)
	uc.Route(apiv1)
	mc.Route(apiv1)
}
