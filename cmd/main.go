package main

import (
	"gIM/internal/Init"
	"gIM/internal/server/bootstrap"
)

func main() {
	Init.Initdb()
	Init.InitLogger()
	bootstrap.Start()
}
