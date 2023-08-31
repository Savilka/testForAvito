package main

import (
	"testForAvito/internal/app"
	"testForAvito/internal/config"
)

var App app.App

func main() {
	cfg := config.MustLoad()

	err := App.Run(cfg)
	if err != nil {
		panic(err)
	}
}
