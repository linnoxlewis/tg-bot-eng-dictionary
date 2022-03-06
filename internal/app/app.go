package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"linnoxlewis/tg-bot-eng-dictionary/internal/api"
	"linnoxlewis/tg-bot-eng-dictionary/internal/config"
	"linnoxlewis/tg-bot-eng-dictionary/internal/db"
	"linnoxlewis/tg-bot-eng-dictionary/internal/dispatcher"
	"linnoxlewis/tg-bot-eng-dictionary/internal/manager"
	"linnoxlewis/tg-bot-eng-dictionary/internal/repository"
	"os"
	"os/signal"
	"time"
)

func Run(ctx context.Context) {
	cfg := config.NewConfig()
	logger := logrus.New()

	database := db.Init(ctx, cfg)
	userRepo := repository.NewUserRepo(database)

	translateService := api.NewSkyEngApi()
	mng := manager.NewManager(logger, userRepo, translateService)

	tgBotDispatcher := dispatcher.NewTgBotDispatcher(cfg, logger, mng)
	go tgBotDispatcher.Run(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

}
