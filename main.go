package main

import (
	"fmt"
	"os"
	"time"
	"watgbridge/database"
	"watgbridge/modules"
	"watgbridge/state"
	"watgbridge/telegram"
	"watgbridge/utils"
	"watgbridge/whatsapp"

	"go.uber.org/zap"
)

func main() {
	// Load configuration file and configs
	cfg := state.State.Config
	cfg.SetDefaults()

	if len(os.Args) > 1 {
		cfg.Path = os.Args[1]
	}

	err := cfg.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config file: %s", err))
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

	// Setup telegram bot
	err = telegram.NewTelegramClient()
	if err != nil {
		logger.Fatal("failed to initialize telegram client",
			zap.Error(err),
		)
		panic(err)
	}
	telegram.AddTelegramHandlers()
	modules.LoadModuleHandlers()

	if cfg.Telegram.RemoveBotCommands {
		err = utils.TgRegisterBotCommands(cfg.Telegram.OwnerID, cfg.Telegram.SkipStartupMessage, state.State.TelegramBot)
		if err != nil {
			logger.Error("failed to set remove commands",
				zap.Error(err),
			)
		}
	} else {
		err = utils.TgRegisterBotCommands(cfg.Telegram.OwnerID, cfg.Telegram.SkipStartupMessage, state.State.TelegramBot, state.State.TelegramCommands...)
		if err != nil {
			logger.Error("failed to set my commands",
				zap.Error(err),
			)
		}
	}
	logger.Sync()

	// Create local location for time
	locLoc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		logger.Fatal("failed to set time zone",
			zap.String("time_zone", cfg.TimeZone),
			zap.Error(err),
		)
		panic(err)
	}
	state.State.LocalLocation = locLoc

	if cfg.FfmpegExecutable == "" && !cfg.Telegram.SkipVideoStickers {
		panic("you need to set your ffmpeg binary location in the config to use video stickers; either skip_video_stickers or provide ffmpeg_executable")
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

	// keep the application running
	state.State.TelegramUpdater.Idle()
}
