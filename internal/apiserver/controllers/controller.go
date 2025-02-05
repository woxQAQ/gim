package controllers

import "github.com/go-fuego/fuego"

type Controller interface {
	Route(sv *fuego.Server)
}
