package state

import (
	"fmt"
	"io"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Path             string `yaml:"-"`
	TimeZone         string `yaml:"time_zone"`
	TimeFormat       string `yaml:"time_format"`
	FfmpegExecutable string `yaml:"ffmpeg_executable"`
	DebugMode        bool   `yaml:"debug_mode"`
	SilentDbLogs     bool   `yaml:"silent_db_logs"`

	Telegram struct {
		BotToken                string  `yaml:"bot_token"`
		ApiUrl                  string  `yaml:"api_url"`
		SudoUsersID             []int64 `yaml:"sudo_users_id"`
		OwnerID                 int64   `yaml:"owner_id"`
		TargetChatID            int64   `yaml:"target_chat_id"`
		SelfHostedAPI           bool    `yaml:"self_hosted_api"`
		SkipVideoStickers       bool    `yaml:"skip_video_stickers"`
		RemoveBotCommands       bool    `yaml:"remove_bot_commands"`
		SendMyPresenceOnReply   bool    `yaml:"send_my_presence_on_reply_on_reply"`
		SendReadReceiptsOnReply bool    `yaml:"send_read_receipts_on_reply"`
		SilentConfirmation      bool    `yaml:"silent_confirmation"`
		ConfirmationType        string  `yaml:"confirmation_type"`
		SkipStartupMessage      bool    `yaml:"skip_startup_message"`
		SpoilerViewOnce         bool    `yaml:"spoiler_as_viewonce"`
		Reactions               bool    `yaml:"reactions"`
	} `yaml:"telegram"`

	WhatsApp struct {
		LoginDatabase struct {
			Type string `yaml:"type"`
			URL  string `yaml:"url"`
		} `yaml:"login_database"`
		StickerMetadata struct {
			PackName   string `yaml:"pack_name"`
			AuthorName string `yaml:"author_name"`
		} `yaml:"sticker_metadata"`
		SessionName                    string   `yaml:"session_name"`
		BrowserName                    string   `yaml:"browser_name"`
		TagAllAllowedGroups            []string `yaml:"tag_all_allowed_groups"`
		IgnoreChats                    []string `yaml:"ignore_chats"`
		StatusIgnoredChats             []string `yaml:"status_ignored_chats"`
		SkipDocuments                  bool     `yaml:"skip_documents"`
		SkipImages                     bool     `yaml:"skip_images"`
		SkipGIFs                       bool     `yaml:"skip_gifs"`
		SkipVideos                     bool     `yaml:"skip_videos"`
		SkipVoiceNotes                 bool     `yaml:"skip_voice_notes"`
		SkipAudios                     bool     `yaml:"skip_audios"`
		SkipStatus                     bool     `yaml:"skip_status"`
		SkipStickers                   bool     `yaml:"skip_stickers"`
		SkipContacts                   bool     `yaml:"skip_contacts"`
		SkipLocations                  bool     `yaml:"skip_locations"`
		SkipProfilePictureUpdates      bool     `yaml:"skip_profile_picture_updates"`
		SkipGroupSettingsUpdates       bool     `yaml:"skip_group_settings_updates"`
		SkipUserAboutUpdates           bool     `yaml:"skip_user_about_updates"`
		SkipChatDetails                bool     `yaml:"skip_chat_details"`
		SkipRevokedMessage             bool     `yaml:"skip_revoked_message"`
		WhatsmeowDebugMode             bool     `yaml:"whatsmeow_debug_mode"`
		SendMyMessagesFromOtherDevices bool     `yaml:"send_my_messages_from_other_devices"`
		CreateThreadForInfoUpdates     bool     `yaml:"create_thread_for_info_updates"`
		AllowEveryoneTagging           bool     `yaml:"allow_everyone_tagging"`
		SkipStartupMessage             bool     `yaml:"skip_startup_message"`
		SkipQrCodeSend                 bool     `yaml:"skip_qr_code"`
		SkipInitialPhotoSend           bool     `yaml:"skip_initial_photo_send"`
	} `yaml:"whatsapp"`

	Database map[string]string `yaml:"database"`
}

func (cfg *Config) LoadConfig() error {
	configFilePath := cfg.Path

	if _, err := os.Stat(configFilePath); err != nil {
		return fmt.Errorf("error with config file path : %s", err)
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return fmt.Errorf("could not open config file : %s", err)
	}
	defer configFile.Close()

	configBody, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("could not read config file : %s", err)
	}

	err = yaml.Unmarshal(configBody, cfg)
	if err != nil {
		return fmt.Errorf("could not parse config file : %s", err)
	}

	deprecatedOptions := GetDeprecatedConfigOptions(cfg)
	if deprecatedOptions != nil {
		fmt.Println("The following options have been deprecated/removed:")
		for num, opt := range deprecatedOptions {
			fmt.Printf("%d. %s: %s\n", num+1, opt.Name, opt.Description)
		}
	}

	return nil
}

func (cfg *Config) SetDefaults() {
	cfg.TimeZone = "UTC"

	cfg.WhatsApp.SessionName = "coco-watg"
	cfg.WhatsApp.BrowserName = "FIREFOX"
	cfg.WhatsApp.LoginDatabase.Type = "sqlite3"
	cfg.WhatsApp.LoginDatabase.URL = "file:coco_wawebstore.db?_foreign_keys=on"
	cfg.WhatsApp.StickerMetadata.PackName = "CocoWaTgBridge"
	cfg.WhatsApp.StickerMetadata.AuthorName = "CocoWaTgBridge"

	cfg.Telegram.ApiUrl = gotgbot.DefaultAPIURL
	cfg.Telegram.ConfirmationType = "emoji"
}
