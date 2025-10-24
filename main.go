package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"watgbridge/database"
	"watgbridge/modules"
	"watgbridge/state"
	"watgbridge/telegram"
	"watgbridge/utils"
	"watgbridge/whatsapp"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

func main() {
	// Load configuration file
	cfg := state.State.Config
	cfg.SetDefaults()

	if len(os.Args) > 1 {
		cfg.Path = os.Args[1]
	}

	err := cfg.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config file: %s", err))
	}

	if cfg.Telegram.APIURL == "" {
		cfg.Telegram.APIURL = gotgbot.DefaultAPIURL
	}

	if cfg.DebugMode {
		developmentConfig := zap.NewDevelopmentConfig()
		developmentConfig.OutputPaths = append(developmentConfig.OutputPaths, "debug.log")
		state.State.Logger, err = developmentConfig.Build()
		if err != nil {
			panic(fmt.Errorf("failed to initialize development logger: %s", err))
		}
		state.State.Logger = state.State.Logger.Named("Coco_WaTgBridge_Dev")
	} else {
		productionConfig := zap.NewProductionConfig()
		state.State.Logger, err = productionConfig.Build()
		if err != nil {
			panic(fmt.Errorf("failed to initialize production logger: %s", err))
		}
		state.State.Logger = state.State.Logger.Named("Coco_WaTgBridge")
	}
	logger := state.State.Logger

	logger.Debug("loaded config file and started logger",
		zap.String("config_path", cfg.Path),
		zap.Bool("development_mode", cfg.DebugMode),
	)
	logger.Sync()

	// Create local location for time
	if cfg.TimeZone == "" {
		cfg.TimeZone = "UTC"
	}
	locLoc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		logger.Fatal("failed to set time zone",
			zap.String("time_zone", cfg.TimeZone),
			zap.Error(err),
		)
	}
	state.State.LocalLocation = locLoc

	if cfg.WhatsApp.SessionName == "" {
		cfg.WhatsApp.SessionName = "coco_watg"
	}

	if cfg.WhatsApp.LoginDatabase.Type == "" || cfg.WhatsApp.LoginDatabase.URL == "" {
		cfg.WhatsApp.LoginDatabase.Type = "sqlite3"
		cfg.WhatsApp.LoginDatabase.URL = "file:coco_wasession.db?_foreign_keys=on"
		logger.Debug("using sqlite3 as WhatsApp login database")
		logger.Sync()
	}

	if cfg.FfmpegExecutable == "" && !cfg.Telegram.SkipVideoStickers {
		ffmpegPath, err := exec.LookPath("ffmpeg")
		if err != nil && !errors.Is(err, exec.ErrDot) {
			logger.Fatal("failed to set ffmpeg executable path",
				zap.Error(err),
			)
			panic("you can't include video stickets without having ffmpeg installed or configured")
		}

		cfg.FfmpegExecutable = ffmpegPath
		logger.Info("setting path to ffmpeg executable",
			zap.String("path", ffmpegPath),
		)
		logger.Sync()

		if err = cfg.SaveConfig(); err != nil {
			logger.Fatal("failed to save config file",
				zap.Error(err),
			)
		}
	}

	// Setup database
	db, err := database.Connect()
	if err != nil {
		logger.Fatal("could not connect to database",
			zap.Error(err),
		)
		panic("could not connect to database")
	}

	state.State.Database = db
	err = database.AutoMigrate()
	if err != nil {
		logger.Fatal("could not migrate database tabels",
			zap.Error(err),
		)
		panic("unable to migrate database")
	}

	err = telegram.NewTelegramClient()
	if err != nil {
		logger.Fatal("failed to initialize telegram client",
			zap.Error(err),
		)
		panic(err)
	}
	err = whatsapp.NewWhatsAppClient()
	if err != nil {
		panic(err)
	}
	logger.Sync()

	state.State.StartTime = time.Now().UTC()

	// s := gocron.NewScheduler(time.UTC)
	// s.TagsUnique()
	// _, _ = s.Every(1).Hour().Tag("foo").Do(func() {
	// 	contacts, err := state.State.WhatsAppClient.Store.Contacts.GetAllContacts(context.Background())
	// 	if err == nil {
	// 		_ = database.ContactNameBulkAddOrUpdate(contacts)
	// 	}
	// })

	state.State.WhatsAppClient.AddEventHandler(whatsapp.WhatsAppEventHandler)
	telegram.AddTelegramHandlers()
	modules.LoadModuleHandlers()

	if cfg.Telegram.RemoveBotCommands {
		err = utils.TgRegisterBotCommands(state.State.TelegramBot)
		if err != nil {
			logger.Error("failed to set my commands to empty",
				zap.Error(err),
			)
		}
	} else {
		err = utils.TgRegisterBotCommands(state.State.TelegramBot, state.State.TelegramCommands...)
		if err != nil {
			logger.Error("failed to set my commands",
				zap.Error(err),
			)
		}
	}

	logger.Sync()

	if !cfg.Telegram.SkipStartupMessage {
		state.State.TelegramBot.SendMessage(cfg.Telegram.OwnerID, "Successfully started Coco_WaTgBridge", &gotgbot.SendMessageOpts{})
	}

	state.State.TelegramUpdater.Idle()
}
