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

/*
type MsgIdPair struct {
	// WhatsApp
	ID            string `gorm:"primaryKey;"` // Message ID
	ParticipantId string // Sender JID
	WaChatId      string // Chat JID

	// Telegram
	TgThreadId int64
	TgMsgId    int64
}

WaMessageId string `gorm:"primaryKey;"` // the actual ID that whatsapp has for the message
	WaSenderJid string // could be jid or lid, or whatever else they decide to add - the datatype can be parsed into a waTypes.ParseJID(string) though
	WaChatJid   string // could be jid or lid, or whatever else they decide to add - the datatype can be parsed into a waTypes.ParseJID(string) though

	// Telegram
	TgThreadId  int64 // thread id (managed by telegram)
	TgMessageId int64 // the actual message id in telegram
*/
