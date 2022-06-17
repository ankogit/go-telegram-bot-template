package main

import (
	"context"
	"errors"
	config "github.com/ankogit/go-telegram-bot-template/configs"
	"github.com/ankogit/go-telegram-bot-template/pkg/server"
	"github.com/ankogit/go-telegram-bot-template/pkg/server/handler"
	"github.com/ankogit/go-telegram-bot-template/pkg/service"
	"github.com/ankogit/go-telegram-bot-template/pkg/storage/stormDB"
	"github.com/ankogit/go-telegram-bot-template/pkg/telegram"
	"github.com/asdine/storm/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Panic(err)
	}

	config := new(config.IniConf)
	config.CheckAndLoadConf("configs" + string(os.PathSeparator) + "config.ini")
	telegramkey := config.GetStringKey("", "telegramkey")

	db, err := storm.Open(config.GetStringKey("", "dbName"))
	defer db.Close()

	chatRepository := stormDB.NewChatRepository(db)

	//Init services
	services := service.NewServices(service.Deps{})

	bot, err := tgbotapi.NewBotAPI(telegramkey)
	if err != nil {
		log.Panic("Wrong key:", telegramkey, err)
	}

	bot.Debug = cfg.Debug

	timeZone, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.New(cron.WithLocation(timeZone))
	defer scheduler.Stop()

	go func() {
		telegramBot := telegram.NewBot(
			bot,
			config,
			cfg.Version,
			cfg.Messages,
			chatRepository,
			services,
		)

		telegramBot.CronInit(scheduler)
		go func() {
			telegramBot.CronStart()
		}()

		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	serverInstance := new(server.Server)
	handlers := handler.NewHandler(services)

	go func() {
		if err := serverInstance.RunHttp(cfg.Port, handlers.InitRoutes()); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	if err := serverInstance.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occurated on shuting down server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatalf("error occurated on db close: %s", err.Error())
	}
	log.Println("Shutdown...")
}
