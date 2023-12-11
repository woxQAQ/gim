package config

import (
	"log"

	"github.com/woxQAQ/gim/internal/global"
	"gopkg.in/ini.v1"
)

func InitConfig() {
	var err error

	global.Cfg, err = ini.Load("internal/config/app.ini")
	if err != nil {
		log.Fatalln("Fail to load config internal/config/app.ini, error:", err)
	}
	initGin()
	initJWT()
	initServer()
	initDB()
}

func initGin() {
	global.RunMode = global.Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func initServer() {
	sec, err := global.Cfg.GetSection("server")
	if err != nil {
		log.Fatalln("Fail to load config section internal/config/app.ini, section server, error:", err)
	}
	global.HttpPort = sec.Key("HTTP_PORT").MustInt(8000)
	global.RecallTimes = sec.Key("RECALL_TIMES").MustInt(100)
}

func initJWT() {
	sec, err := global.Cfg.GetSection("app")
	if err != nil {
		log.Fatalln("Fail to load config section internal/config/app.ini, section app, error:", err)
	}
	global.JwtSecret = sec.Key("JWT_SECRET").String()
	if global.JwtSecret == "" {
		log.Fatalln("You forget to set JWT_SECRET option, failed to init server")
	}
}

func initDB() {
	sec, err := global.Cfg.GetSection("database")
	if err != nil {
		log.Fatalln("Fail to load config section internal/config/app.ini, section database, error:", err)
	}
	global.DBName = sec.Key("NAME").String()
	global.DBHost = sec.Key("HOST").String()
	global.DBPwd = sec.Key("PASSWORD").String()
	global.DBUser = sec.Key("USER").String()

}
