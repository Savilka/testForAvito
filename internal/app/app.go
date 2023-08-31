package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testForAvito/internal/api/v1"
	"testForAvito/internal/config"
	"testForAvito/internal/storage/postgres"
	"testForAvito/internal/utils"
)

type App struct {
	Db     *postgres.Storage
	Router *gin.Engine
}

func (a *App) Run(cfg *config.Config) error {
	const op = "app.Run"

	var err error
	a.Db, err = postgres.Connect(cfg.DataBaseUrl)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		err := utils.CheckTtlForSegments(a.Db)
		if err != nil {
			panic(err)
		}
	}()

	a.Router = api.InitRouter(a.Db)

	err = a.Router.Run(cfg.Address + ":" + cfg.Port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
