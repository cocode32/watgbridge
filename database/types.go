package database

import (
	"errors"
	"watgbridge/state"
)

type CocoContact struct {
	ID           int32 `gorm:"primaryKey;autoIncrement"`
	Lid          string
	Jid          string
	Name         string
	FullName     string
	PushName     string
	BusinessName string
}

type CocoChatThread struct {
	ID            int32 `gorm:"primaryKey;autoIncrement"`
	CocoContactId int32
	ThreadId      int64
}

type MsgIdPair struct {
	// WhatsApp
	ID            string `gorm:"primaryKey;"` // Message ID
	ParticipantId string // Sender JID
	WaChatId      string // Chat JID

	// Telegram
	TgThreadId int64
	TgMsgId    int64
}

type ChatEphemeralSettings struct {
	ID             string `gorm:"primaryKey;"` // WhatsApp Chat ID
	IsEphemeral    bool
	EphemeralTimer uint32
}

func AutoMigrate() error {
	db := state.State.Database
	autoMigrateError := db.AutoMigrate(
		&MsgIdPair{},
		&ChatEphemeralSettings{},
		&CocoContact{},
		&CocoChatThread{},
	)

	migrateError := MigrateDatabase(db)
	return errors.Join(autoMigrateError, migrateError)
}
