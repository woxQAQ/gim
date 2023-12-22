package config

import (
	"log"

	"gopkg.in/ini.v1"
)

func InitConfig() {
	var err error

	Cfg, err = ini.Load("config/app.ini")
	if err != nil {
		log.Fatalln("Fail to load config internal/config/app.ini, error:", err)
	}
	initGin()
	initJWT()
	initServer()
	initDB()
}

func initGin() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func initServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalln("Fail to load config section internal/config/app.ini, section server, error:", err)
	}
	HttpPort = sec.Key("HTTP_PORT").MustInt(8000)
	RecallTimes = sec.Key("RECALL_TIMES").MustInt(100)
}

func initJWT() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalln("Fail to load config section internal/config/app.ini, section app, error:", err)
	}
	JwtSecret = sec.Key("JWT_SECRET").String()
	if JwtSecret == "" {
		log.Fatalln("You forget to set JWT_SECRET option, failed to init server")
	}
}

func initDB() {
	sec, err := Cfg.GetSection("database")
	if err != nil {
		log.Fatalln("Fail to load config section internal/config/app.ini, section database, error:", err)
	}
	DBName = sec.Key("NAME").String()
	DBHost = sec.Key("HOST").String()
	DBPwd = sec.Key("PASSWORD").String()
	DBUser = sec.Key("USER").String()

}
